package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"

	// "github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	// "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
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
	ServerURL          string `koanf:"url" validate:"required,url"`
	ReadTimeout        int    `koanf:"read_timeout" validate:"required,min=1"`
	WriteTimeout       int    `koanf:"write_timeout" validate:"required,min=1"`
	IdleTimeout        int    `koanf:"idle_timeout" validate:"required,min=1"`
	CORSAllowedOrigins string `koanf:"cors_allowed_origins" validate:"required"`
}

// DatabaseConfig contains PostgreSQL database configuration
type DatabaseConfig struct {
	DatabaseURL     string `koanf:"url" validate:"required"`
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required,min=1,max=65535"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password" validate:"required"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required,oneof=disable allow prefer require verify-ca verify-full"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required,min=1"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required,min=0"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required,min=1"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required,min=0"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `koanf:"host" validate:"required"`
	Port     int    `koanf:"port" validate:"required,min=1,max=65535"`
	Username string `koanf:"username"`
	Password string `koanf:"password" validate:"required"`
	Address  string `koanf:"address" validate:"required"`
}

// RabbitMQConfig contains RabbitMQ message queue configuration
type RabbitMQConfig struct {
	URL        string `koanf:"url" validate:"required"`
	PrivateURL string `koanf:"private_url" validate:"required"`
	NodeName   string `koanf:"node_name" validate:"required"`
	User       string `koanf:"user" validate:"required"`
	Password   string `koanf:"password" validate:"required"`
}

// EmailConfig contains email service configuration
type EmailConfig struct {
	ResendKey string `koanf:"resend_key" validate:"required"`
}

// S3Config contains AWS S3 backup configuration
type S3Config struct {
	BackupEnabled bool   `koanf:"backup_enabled"`
	AccessKey     string `koanf:"access_key" validate:"required"`
	SecretKey     string `koanf:"secret_key" validate:"required"`
	Region        string `koanf:"region" validate:"required"`
	Bucket        string `koanf:"bucket" validate:"required"`
	Prefix        string `koanf:"prefix"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
}

// LoadConfig loads and validates the configuration from environment variables and .env file
func LoadConfig(envFilePath string) (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("BOILERPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BOILERPLATE_"))
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}
	// Validate the config
	if err := ValidateConfig(mainConfig); err != nil {
		return nil, err
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

func ValidateConfig(cfg *Config) error {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		var errMessages []string
		for _, validationErr := range err.(validator.ValidationErrors) {
			errMessages = append(errMessages, formatValidationError(validationErr))
		}
		return fmt.Errorf("config validation failed:\n  - %s",
			strings.Join(errMessages, "\n  - "))
	}

	return nil
}

// formatValidationError formats a validation error for better readability
func formatValidationError(err validator.FieldError) string {
	return fmt.Sprintf(
		"field '%s' failed validation '%s' (value: %v)",
		err.StructNamespace(),
		err.Tag(),
		err.Value(),
	)
}
