package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config is the main application configuration struct
type Config struct {
	Primary       PrimaryConfig        `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	RabbitMQ      RabbitMQConfig       `koanf:"rabbitmq" validate:"required"`
	Email         EmailConfig          `koanf:"email" validate:"required"`
	S3            S3Config             `koanf:"s3" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
}

// PrimaryConfig contains basic environment configuration
type PrimaryConfig struct {
	Env string `koanf:"env" validate:"required,oneof=development staging production"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port               string `koanf:"port" validate:"required,numeric"`
	ServerURL          string `koanf:"server_url" validate:"required,url"`
	ReadTimeout        int    `koanf:"read_timeout" validate:"required,min=1"`
	WriteTimeout       int    `koanf:"write_timeout" validate:"required,min=1"`
	IdleTimeout        int    `koanf:"idle_timeout" validate:"required,min=1"`
	CORSAllowedOrigins string `koanf:"cors_allowed_origins" validate:"required"`
}

// DatabaseConfig contains PostgreSQL database configuration
type DatabaseConfig struct {
	DatabaseURL     string `koanf:"database_url" validate:"required"`
	Host            string `koanf:"database_host" validate:"required"`
	Port            int    `koanf:"database_port" validate:"required,min=1,max=65535"`
	User            string `koanf:"database_user" validate:"required"`
	Password        string `koanf:"database_password" validate:"required"`
	Name            string `koanf:"database_name" validate:"required"`
	SSLMode         string `koanf:"database_ssl_mode" validate:"required,oneof=disable allow prefer require verify-ca verify-full"`
	MaxOpenConns    int    `koanf:"database_max_open_conns" validate:"required,min=1"`
	MaxIdleConns    int    `koanf:"database_max_idle_conns" validate:"required,min=0"`
	ConnMaxLifetime int    `koanf:"database_conn_max_lifetime" validate:"required,min=1"`
	ConnMaxIdleTime int    `koanf:"database_conn_max_idle_time" validate:"required,min=0"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `koanf:"redis_host" validate:"required"`
	Port     int    `koanf:"redis_port" validate:"required,min=1,max=65535"`
	Username string `koanf:"redis_username"`
	Password string `koanf:"redis_password" validate:"required"`
}

// RabbitMQConfig contains RabbitMQ message queue configuration
type RabbitMQConfig struct {
	URL        string `koanf:"rabbitmq_url" validate:"required"`
	PrivateURL string `koanf:"rabbitmq_private_url" validate:"required"`
	NodeName   string `koanf:"rabbitmq_nodename" validate:"required"`
	User       string `koanf:"rabbitmq_default_user" validate:"required"`
	Password   string `koanf:"rabbitmq_default_pass" validate:"required"`
}

// EmailConfig contains email service configuration
type EmailConfig struct {
	ResendKey string `koanf:"resend_key" validate:"required"`
}

// S3Config contains AWS S3 backup configuration
type S3Config struct {
	BackupEnabled   bool   `koanf:"s3_backup_enabled"`
	BucketName      string `koanf:"bucket_name" validate:"required_if=S3BackupEnabled true"`
	BucketRegion    string `koanf:"bucket_region" validate:"required_if=S3BackupEnabled true"`
	AccessKey       string `koanf:"access_key" validate:"required_if=S3BackupEnabled true"`
	SecretAccessKey string `koanf:"secret_access_key" validate:"required_if=S3BackupEnabled true"`
	Prefix          string `koanf:"s3_prefix"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	SecretKey    string `koanf:"secret_key" validate:"required"`
	ResendAPIKey string `koanf:"resend_api_key" validate:"required"`
}

// LoadConfig loads and validates the configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	// Load from .env file if it exists
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil && !strings.Contains(err.Error(), "no such file") {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Load from environment variables (overrides .env)
	if err := k.Load(env.Provider("", ".", nil), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Unmarshal into Config struct
	mainConfig := &Config{}
	err := k.Unmarshal("", mainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the configuration
	validate := validator.New()
	err = validate.Struct(mainConfig)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	// Set default observability config if not provided
	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	// Override service name and environment from primary config
	mainConfig.Observability.ServiceName = "Fortress_API"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	// Validate observability config
	if err := mainConfig.Observability.Validate(); err != nil {
		// logger.Fatal().Err(err).Msg("invalid observability config")
	}
	return mainConfig, nil
}
