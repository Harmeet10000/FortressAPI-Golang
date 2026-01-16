package workers

import (
	"log"

	"github.com/Harmeet10000/Fortress_API/src/internal/config"
	"github.com/Harmeet10000/Fortress_API/src/internal/connections"
	"github.com/Harmeet10000/Fortress_API/src/internal/features/auth"
	"github.com/hibiken/asynq"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv := connections.NewAsynqServer(cfg)

	mux := asynq.NewServeMux()

	// Register feature-based handlers
	mux.HandleFunc(auth.TypeWelcomeEmail, auth.HandleWelcomeEmailTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

