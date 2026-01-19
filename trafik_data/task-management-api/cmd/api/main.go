package main

import (
	"log"
	"os"

	_ "github.com/yourusername/task-management-api/docs" // Swagger docs
	
	"github.com/yourusername/task-management-api/internal/app"
	"github.com/yourusername/task-management-api/internal/config"
	"github.com/yourusername/task-management-api/internal/feature/category"
	"github.com/yourusername/task-management-api/internal/feature/comment"
	"github.com/yourusername/task-management-api/internal/feature/todo"
	"github.com/yourusername/task-management-api/internal/logger"
)

// @title Task Management API
// @version 1.0
// @description A professional task management API built with Go
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	loggerService := logger.New(&cfg.Log)
	logger := loggerService.GetLogger()

	logger.Info().Msg("Starting Task Management API")

	// Create server instance
	srv := server.NewServer(cfg, loggerService)

	// Initialize connections
	if err := srv.InitializeConnections(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize connections")
	}

	// Initialize repositories
	categoryRepo := category.NewRepository(srv.DB, logger)
	todoRepo := todo.NewRepository(srv.DB, logger)
	commentRepo := comment.NewRepository(srv.DB, logger)

	// Initialize services
	categoryService := category.NewService(categoryRepo, logger)
	todoService := todo.NewService(todoRepo, categoryRepo, logger)
	commentService := comment.NewService(commentRepo, todoRepo, logger)

	// Initialize handlers
	categoryHandler := category.NewHandler(categoryService, logger)
	todoHandler := todo.NewHandler(todoService, logger)
	commentHandler := comment.NewHandler(commentService, logger)

	// Setup HTTP server
	srv.SetupHTTPServer(categoryHandler, todoHandler, commentHandler)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	logger.Info().
		Str("address", cfg.GetServerAddr()).
		Msg("Server started successfully")

	// Wait for shutdown signal
	srv.WaitForShutdown()
}
