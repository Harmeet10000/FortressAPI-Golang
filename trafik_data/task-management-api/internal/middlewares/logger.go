package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// RequestLogger creates a request logging middleware
func RequestLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Process request
			err := next(c)

			// Log request
			duration := time.Since(start)
			
			logEvent := logger.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("path", c.Path()).
				Int("status", res.Status).
				Dur("duration", duration).
				Str("ip", c.RealIP()).
				Str("user_agent", req.UserAgent())

			if err != nil {
				logEvent = logEvent.Err(err)
			}

			logEvent.Msg("HTTP request")

			return err
		}
	}
}
