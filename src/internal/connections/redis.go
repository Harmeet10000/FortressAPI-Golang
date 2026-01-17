package connections

import (
	"context"
	"fmt"
	"time"

	"github.com/Harmeet10000/Fortress_API/src/internal/config"
	"github.com/redis/go-redis/v9"
)

// Configuration options:
// - maxRetriesPerRequest: 3 (max retry attempts for failed commands)
// - retryDelayOnFailover: 100ms (delay between retries on failover)
// - keepAlive: 120 seconds (TCP keep-alive interval)
// - family: 4 (IPv4 only)
// - db: 0 (default database)
// - connectTimeout: 120 seconds (connection establishment timeout)
// - commandTimeout: 5 seconds (individual command timeout)
// - enableAutoPipelining: true (automatic command pipelining)
func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		// Connection settings
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:       0, // default DB
		Protocol: 3, // RESP 3

		// Authentication
		Username: cfg.Username,
		Password: cfg.Password,

		// Timeouts
		DialTimeout:  120 * time.Second,  // connectTimeout in JS
		ReadTimeout:  5 * time.Second,    // commandTimeout in JS
		WriteTimeout: 5 * time.Second,    // commandTimeout in JS

		// Connection pool settings
		MaxRetries:         3,                  // maxRetriesPerRequest
		// MaxConnAge:         0,                  // connection max age (no limit)
		PoolSize:           10,                 // default: 10 connections per CPU
		MinIdleConns:       5,                  // minimum idle connections
		MaxIdleConns:       10,                 // maximum idle connections
		ConnMaxIdleTime:    5 * time.Minute,    // close idle connections after 5 min
		ConnMaxLifetime:    0,                  // no max lifetime

		// Automatic failover and recovery
		// RetryOnTimeout:         true,                // retry commands on timeout
		MinRetryBackoff:        8 * time.Millisecond,
		MaxRetryBackoff:        512 * time.Millisecond,
		ContextTimeoutEnabled:  true,

		// Other settings
		PoolFIFO:            false,
		// DisableIdentifier:  false,
		// IdentifierTemplate:  "redis-go-{lib-version}",
	})

	return client
}

// NewRedisClientWithConfig creates a Redis client using configuration struct
// This version allows for flexible configuration from environment variables or config files
func NewRedisClientWithConfig(cfg *RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		// Connection settings
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:       cfg.DB,
		Protocol: 3, // RESP 3

		// Authentication
		Username: cfg.Username,
		Password: cfg.Password,

		// Timeouts (convert milliseconds to duration)
		DialTimeout:  time.Duration(cfg.ConnectTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.CommandTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.CommandTimeout) * time.Millisecond,

		// Connection pool settings
		MaxRetries:        cfg.MaxRetries,
		// MaxConnAge:        0,
		PoolSize:          cfg.PoolSize,
		MinIdleConns:      cfg.MinIdleConns,
		MaxIdleConns:      cfg.MaxIdleConns,
		ConnMaxIdleTime:   time.Duration(cfg.ConnMaxIdleTime) * time.Second,
		ConnMaxLifetime:   0,

		// Automatic failover and recovery
		// RetryOnTimeout:        cfg.RetryOnTimeout,
		MinRetryBackoff:       8 * time.Millisecond,
		MaxRetryBackoff:       512 * time.Millisecond,
		ContextTimeoutEnabled: cfg.ContextTimeoutEnabled,

		// Other settings
		PoolFIFO:            false,
		// DisableIndentifier:  false,
		// IdentifierTemplate:  "redis-go-{lib-version}",
	})

	return client
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	DB                    int
	MaxRetries            int
	ConnectTimeout        int // milliseconds
	CommandTimeout        int // milliseconds
	PoolSize              int
	MinIdleConns          int
	MaxIdleConns          int
	ConnMaxIdleTime       int  // seconds
	RetryOnTimeout        bool
	ContextTimeoutEnabled bool
}

// PingRedis checks Redis connection health
func PingRedis(ctx context.Context, client *redis.Client) error {
	result := client.Ping(ctx)
	return result.Err()
}

// CloseRedis gracefully closes the Redis client connection
func CloseRedis(client *redis.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}

