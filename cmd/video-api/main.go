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
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/trieuvy/video-ranking/internal/database"
	"github.com/trieuvy/video-ranking/internal/handlers"
	"github.com/trieuvy/video-ranking/internal/models"
	"github.com/trieuvy/video-ranking/internal/redis"
	"github.com/trieuvy/video-ranking/internal/repositories"
	"github.com/trieuvy/video-ranking/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Initialize Redis client
	redis, err := redis.ConnectRedis()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	// Initialize MySQL database connection
	db,err := database.ConnectDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	// Migrate database schema
	if err := db.AutoMigrate(&models.Video{}, &models.User{}, &models.Interaction{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
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
	interactionHandler := handlers.NewInteractionHandler(interactionService, videoService)

	// Initialize router
	r := mux.NewRouter()

	// Register routes
	videoHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	interactionHandler.RegisterRoutes(r)

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

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
