package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type Hub struct {
	clients map[*Client]bool
	lock    sync.RWMutex
}

var hub = &Hub{
	clients: make(map[*Client]bool),
}

func RegisterClient(conn *websocket.Conn) *Client {
	client := &Client{Conn: conn}
	hub.lock.Lock()
	hub.clients[client] = true
	hub.lock.Unlock()
	return client
}

func UnregisterClient(client *Client) {
	hub.lock.Lock()
	delete(hub.clients, client)
	hub.lock.Unlock()
}

func Broadcast(message []byte) {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	for client := range hub.clients {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}
	client := RegisterClient(conn)
	defer UnregisterClient(client)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		Broadcast(message)
	}
}
