package connections

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"github.com/Harmeet10000/Fortress_API/src/internal/config"
)

func NewAsynqClient(cfg *config.Config) *asynq.Client {
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:         redisAddr,
		Username:     cfg.Redis.Username,
		Password:     cfg.Redis.Password,
		DB:           0,                 // Equivalent to db: 0
		DialTimeout:  120 * time.Second, // Equivalent to connectTimeout
		ReadTimeout:  5 * time.Second,   // Command timeout
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		// MinIdleConns: 5,
	})
}

func NewAsynqServer(cfg *config.Config) *asynq.Server {
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:         redisAddr,
			Username:     cfg.Redis.Username,
			Password:     cfg.Redis.Password,
			DB:           0,
			DialTimeout:  120 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			PoolSize:     10,
			// MinIdleConns: 5,
		},
		asynq.Config{
			Concurrency:     10,
			Queues:          map[string]int{"critical": 6, "default": 3, "low": 1},
			RetryDelayFunc:  defaultRetryDelay,  // Custom retry logic
			// ErrorHandler:    customErrorHandler, // Handle failures
			ShutdownTimeout: 30 * time.Second,
		},
	)
}

// Custom retry delay for tasks (not connections)
func defaultRetryDelay(retryCount int, err error, task *asynq.Task) time.Duration {
	return time.Duration(retryCount) * time.Second
}
