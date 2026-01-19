package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/Harmeet10000/Fortress_API/src/internal/app"
	"github.com/Harmeet10000/Fortress_API/src/internal/config"
	"github.com/Harmeet10000/Fortress_API/src/internal/handler"
	"github.com/Harmeet10000/Fortress_API/src/internal/logger"
	"github.com/Harmeet10000/Fortress_API/src/internal/repository"
	"github.com/Harmeet10000/Fortress_API/src/internal/router"
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

	// if cfg.Primary.Env != "local" {
	// 	if err := connections.Migrate(context.Background(), &log, cfg); err != nil {
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

	// Start server
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
