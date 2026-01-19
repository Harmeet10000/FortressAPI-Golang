package connections

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/config"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig, logger *zerolog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	logger.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Int("db", cfg.DB).
		Msg("Redis connection established")

	return client, nil
}

// CloseRedis closes the Redis connection
func CloseRedis(client *redis.Client, logger *zerolog.Logger) {
	if client != nil {
		if err := client.Close(); err != nil {
			logger.Error().Err(err).Msg("Error closing Redis connection")
		} else {
			logger.Info().Msg("Redis connection closed")
		}
	}
}
