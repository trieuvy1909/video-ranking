package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/trieuvy/video-ranking/configs/database"
	"github.com/trieuvy/video-ranking/configs/redis"
	_ "github.com/trieuvy/video-ranking/docs"
	"github.com/trieuvy/video-ranking/internal/handlers"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/trieuvy/video-ranking/internal/repositories"
	"github.com/trieuvy/video-ranking/internal/services"
	"github.com/trieuvy/video-ranking/internal/ws"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Environment variables loaded successfully")
	// Initialize Redis client
	redis, err := redis.ConnectRedis()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return
	}
	log.Println("Redis client initialized successfully")
	// Initialize MySQL database connection
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	log.Println("Database connection established successfully")
	// Migrate database schema
	if err := db.AutoMigrate(&models.Video{}, &models.User{}, &models.Interaction{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return
	}
	log.Println("Database migration completed successfully")
	// Initialize repositories
	videoRepo := repositories.NewVideoRepository(db)
	userRepo := repositories.NewUserRepository(db)
	interactionRepo := repositories.NewInteractionRepository(db)

	// Initialize services
	videoService := services.NewVideoService(videoRepo, redis)
	userService := services.NewUserService(userRepo)
	interactionService := services.NewInteractionService(interactionRepo)

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(videoService)
	userHandler := handlers.NewUserHandler(userService)
	interactionHandler := handlers.NewInteractionHandler(interactionService, videoService, userService)

	// Initialize router
	r := mux.NewRouter()

	// Register routes
	videoHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	interactionHandler.RegisterRoutes(r)

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(r),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		_ = browser.OpenURL("http://localhost:8080/swagger/index.html")
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	r.HandleFunc("/ws", ws.WsHandler).Methods("GET")
	<-sigChan
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped successfully")
}
