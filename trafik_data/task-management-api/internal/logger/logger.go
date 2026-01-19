package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/config"
)

// LoggerService provides structured logging capabilities
type LoggerService struct {
	logger *zerolog.Logger
}

// New creates a new LoggerService instance
func New(cfg *config.LogConfig) *LoggerService {
	var output io.Writer = os.Stdout

	// Configure based on format
	if cfg.Format == "console" || cfg.Format == "pretty" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Set log level
	level := parseLogLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return &LoggerService{
		logger: &logger,
	}
}

// GetLogger returns the underlying zerolog.Logger
func (s *LoggerService) GetLogger() *zerolog.Logger {
	return s.logger
}

// Debug logs a debug message
func (s *LoggerService) Debug(msg string, fields map[string]interface{}) {
	event := s.logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Info logs an info message
func (s *LoggerService) Info(msg string, fields map[string]interface{}) {
	event := s.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (s *LoggerService) Warn(msg string, fields map[string]interface{}) {
	event := s.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs an error message
func (s *LoggerService) Error(msg string, err error, fields map[string]interface{}) {
	event := s.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (s *LoggerService) Fatal(msg string, err error, fields map[string]interface{}) {
	event := s.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// WithContext returns a new logger with additional context
func (s *LoggerService) WithContext(fields map[string]interface{}) *LoggerService {
	ctx := s.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	newLogger := ctx.Logger()
	return &LoggerService{logger: &newLogger}
}

// parseLogLevel parses string log level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
