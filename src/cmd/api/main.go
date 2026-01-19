package main

import (
	"context"
	// "fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"
	"errors"
	"net/http"

	"github.com/Harmeet10000/Fortress_API/src/internal/app"
	"github.com/Harmeet10000/Fortress_API/src/internal/config"
	"github.com/Harmeet10000/Fortress_API/src/internal/handler"
	"github.com/Harmeet10000/Fortress_API/src/internal/logger"
	"github.com/Harmeet10000/Fortress_API/src/internal/repository"
	router "github.com/Harmeet10000/Fortress_API/src/internal/router"
	"github.com/Harmeet10000/Fortress_API/src/internal/service"
)

const DefaultContextTimeout = 30

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	envPath := filepath.Join(wd, ".env")
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize New Relic logger service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)
	log.Info().Msg("server loaded successfully")

	// if cfg.Primary.Env != "local" {
	// 	if err := .Migrate(context.Background(), &log, cfg); err != nil {
	// 		log.Fatal().Err(err).Msg("failed to migrate database")
	// 	}
	// }

	// Initialize server
	srv, err := app.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Initialize repositories, services, and handlers
	repos := repository.NewRepositories(srv)
	services, serviceErr := service.NewServices(srv, repos)
	if serviceErr != nil {
		log.Fatal().Err(serviceErr).Msg("could not create services")
	}
	handlers := handler.NewHandlers(srv, services)

	// Initialize router
	r := router.NewRouter(srv, handlers, services)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// // Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	log.Info().Msg("server exited properly")
}
// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"your-project/internal/connections"
// 	"your-project/internal/middleware"

// 	// Features
// 	"your-project/internal/features/auth"
// 	"your-project/internal/features/users"
// 	"your-project/internal/features/files"
// 	"your-project/internal/features/payments"

// 	"github.com/labstack/echo/v4"
// 	echomiddleware "github.com/labstack/echo/v4/middleware"
// )

// func main() {
// 	// 1. Initialize Infrastructure (Database, Redis, etc.)
// 	db := connections.NewPostgresDB()
// 	asynqClient := connections.NewAsynqClient() // For features that need to enqueue tasks

// 	// 2. Initialize Echo
// 	e := echo.New()

// 	// 3. Global Middlewares
// 	e.Use(echomiddleware.Logger())
// 	e.Use(echomiddleware.Recover())
// 	e.Use(echomiddleware.CORS())
// 	e.Use(middleware.CustomAuthMiddleware) // Your own internal middleware

// 	// 4. Feature Initialization & Route Registration
// 	// We pass the DB and Echo instance to each feature
// 	apiGroup := e.Group("/api/v1")

// 	auth.RegisterHandlers(apiGroup, db, asynqClient)
// 	users.RegisterHandlers(apiGroup, db)
// 	files.RegisterHandlers(apiGroup, db)
// 	payments.RegisterHandlers(apiGroup, db)

// 	// 5. Start Server with Graceful Shutdown
// 	go func() {
// 		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
// 			e.Logger.Fatal("shutting down the server")
// 		}
// 	}()

// 	// Wait for interrupt signal to gracefully shutdown the server
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt)
// 	<-quit

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := e.Shutdown(ctx); err != nil {
// 		e.Logger.Fatal(err)
// 	}
// }
