package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
	
	"github.com/yourusername/task-management-api/internal/config"
	"github.com/yourusername/task-management-api/internal/connections"
	"github.com/yourusername/task-management-api/internal/feature/category"
	"github.com/yourusername/task-management-api/internal/feature/comment"
	"github.com/yourusername/task-management-api/internal/feature/todo"
	"github.com/yourusername/task-management-api/internal/logger"
	"github.com/yourusername/task-management-api/internal/middlewares"
)

// Server represents the application server
type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *logger.LoggerService
	DB            *connections.Database
	Redis         *redis.Client
	httpServer    *http.Server
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, loggerService *logger.LoggerService) *Server {
	return &Server{
		Config:        cfg,
		Logger:        loggerService.GetLogger(),
		LoggerService: loggerService,
	}
}

// InitializeConnections initializes database and Redis connections
func (s *Server) InitializeConnections() error {
	// Initialize database connection
	db, err := connections.NewDatabase(&s.Config.Database, s.Logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.DB = db
	s.Logger.Info().Msg("Database connection initialized")

	// Initialize Redis connection
	redisClient, err := connections.NewRedisClient(&s.Config.Redis, s.Logger)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	s.Redis = redisClient
	s.Logger.Info().Msg("Redis connection initialized")

	return nil
}

// SetupHTTPServer sets up the HTTP server with all routes and middlewares
func (s *Server) SetupHTTPServer(
	categoryHandler *category.Handler,
	todoHandler *todo.Handler,
	commentHandler *comment.Handler,
) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup middlewares
	middlewares.Setup(e, s.Logger)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		}

		// Check database
		if err := s.DB.Health(ctx); err != nil {
			health["database"] = "unhealthy"
			health["status"] = "degraded"
		} else {
			health["database"] = "healthy"
		}

		// Check Redis
		if err := s.Redis.Ping(ctx).Err(); err != nil {
			health["redis"] = "unhealthy"
			health["status"] = "degraded"
		} else {
			health["redis"] = "healthy"
		}

		return c.JSON(http.StatusOK, health)
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register feature routes
	category.RegisterRoutes(e, categoryHandler)
	todo.RegisterRoutes(e, todoHandler)
	comment.RegisterRoutes(e, commentHandler)

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.Config.GetServerAddr(),
		Handler:      e,
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
	}

	s.Logger.Info().
		Str("address", s.Config.GetServerAddr()).
		Msg("HTTP server configured")
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.Logger.Info().
		Str("address", s.httpServer.Addr).
		Msg("Starting HTTP server")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.Logger.Info().Msg("Initiating graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), s.Config.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if s.httpServer != nil {
		s.Logger.Info().Msg("Shutting down HTTP server")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.Logger.Error().Err(err).Msg("Error shutting down HTTP server")
		}
	}

	// Close Redis connection
	if s.Redis != nil {
		connections.CloseRedis(s.Redis, s.Logger)
	}

	// Close database connection
	if s.DB != nil {
		s.DB.Close()
	}

	s.Logger.Info().Msg("Graceful shutdown completed")
	return nil
}

// WaitForShutdown blocks until a shutdown signal is received
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	s.Logger.Info().Msg("Shutdown signal received")

	if err := s.Shutdown(); err != nil {
		s.Logger.Fatal().Err(err).Msg("Failed to shutdown gracefully")
	}
}
