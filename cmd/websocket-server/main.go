package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gorilla/websocket"
    "github.com/joho/godotenv"
    "github.com/trieuvy/video-ranking/internal/ws"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins in development
    },
}

type Hub struct {
    clients    map[*websocket.Conn]bool
    broadcast  chan []byte
    register   chan *websocket.Conn
    unregister chan *websocket.Conn
    mutex      sync.Mutex
}

func newHub() *Hub {
    return &Hub{
        clients:    make(map[*websocket.Conn]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *websocket.Conn),
        unregister: make(chan *websocket.Conn),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            log.Printf("New client connected. Total clients: %d", len(h.clients))
            h.mutex.Unlock()
        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                client.Close()
                log.Printf("Client disconnected. Total clients: %d", len(h.clients))
            }
            h.mutex.Unlock()
        case message := <-h.broadcast:
            h.mutex.Lock()
            for client := range h.clients {
                err := client.WriteMessage(websocket.TextMessage, message)
                if err != nil {
                    log.Printf("Error broadcasting to client: %v", err)
                    client.Close()
                    delete(h.clients, client)
                }
            }
            h.mutex.Unlock()
        }
    }
}

func (h *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }

    h.register <- conn

    // Handle incoming messages (if needed)
    go func() {
        defer func() {
            h.unregister <- conn
            conn.Close()
        }()

        for {
            _, _, err := conn.ReadMessage()
            if err != nil {
                if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    log.Printf("WebSocket error: %v", err)
                }
                break
            }
        }
    }()
}

func (h *Hub) broadcastUpdate(update ws.Update) {
    message, err := json.Marshal(update)
    if err != nil {
        log.Printf("Error marshaling update: %v", err)
        return
    }
    log.Printf("Broadcasting update: %s", string(message))
    h.broadcast <- message
}

func main() {
    // Đọc file .env nếu có
    if _, err := os.Stat(".env"); err == nil {
        err := godotenv.Load()
        if err != nil {
            log.Printf("Error loading .env file: %v", err)
        } else {
            log.Println("Loaded environment variables from .env file")
        }
    }

    // Create new hub
    hub := newHub()

    // Start hub
    go hub.run()

    // Kết nối Redis để nhận thông báo trending videos
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        redisURL = "localhost:6379" // Default fallback
        log.Printf("REDIS_URL not set, using default: %s", redisURL)
    }

    redisClient := redis.NewClient(&redis.Options{
        Addr: redisURL,
    })

    // Test kết nối Redis
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        log.Printf("Warning: Failed to connect to Redis: %v", err)
    } else {
        log.Println("Connected to Redis successfully")

        // Đăng ký lắng nghe kênh ws:updates
        pubsub := redisClient.Subscribe(context.Background(), "ws:updates")
        defer pubsub.Close()

        // Xử lý thông báo từ Redis trong goroutine riêng
        go func() {
            log.Printf("Listening for trending videos updates via Redis...")
            for {
                // Nhận thông báo
                msg, err := pubsub.ReceiveMessage(context.Background())
                if err != nil {
                    log.Printf("Error receiving Redis message: %v", err)
                    time.Sleep(time.Second) // Tạm dừng trước khi thử lại
                    continue
                }

                log.Printf("Received Redis message: %s", msg.Payload)

                // Parse thông báo
                var update map[string]interface{}
                if err := json.Unmarshal([]byte(msg.Payload), &update); err != nil {
                    log.Printf("Error parsing Redis message: %v", err)
                    continue
                }

                // Tạo update để gửi qua websocket
                wsUpdate := ws.Update{
                    Type:    ws.UpdateType(update["type"].(string)),
                    Payload: update,
                }

                // Broadcast tới tất cả clients
                hub.broadcastUpdate(wsUpdate)
                log.Printf("Broadcast trending videos to %d clients", len(hub.clients))
            }
        }()
    }

    // Create router
    mux := http.NewServeMux()

    // WebSocket endpoint
    mux.HandleFunc("/ws", hub.handleWebSocket)

    // Add a simple health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("WebSocket server is running"))
    })

    // Create server
    server := &http.Server{
        Addr:    ":8081",
        Handler: mux,
    }

    // Start server
    go func() {
        log.Printf("WebSocket server is running on port 8081")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan
    log.Println("Shutting down server...")

    // Create shutdown context with timeout
    ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Shutdown server
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }

    log.Println("Server stopped successfully")
}