package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// Setup configures all application middlewares
func Setup(e *echo.Echo, logger *zerolog.Logger) {
	// Custom error handler
	e.HTTPErrorHandler = ErrorHandler(logger)

	// Recover from panics
	e.Use(middleware.Recover())

	// Request ID
	e.Use(middleware.RequestID())

	// Custom request logger
	e.Use(RequestLogger(logger))

	// CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Configure based on your needs
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// Gzip compression
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Security headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// Body limit
	e.Use(middleware.BodyLimit("10M"))

	// Rate limiting (optional - configure as needed)
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
}
