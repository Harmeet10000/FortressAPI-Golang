package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/errs"
)

// ErrorHandler is a custom error handler for Echo
func ErrorHandler(logger *zerolog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// Don't process if response already committed
		if c.Response().Committed {
			return
		}

		var (
			code    int
			message interface{}
		)

		// Check if it's our custom AppError
		if appErr, ok := errs.IsAppError(err); ok {
			code = appErr.StatusCode
			message = appErr.ToErrorResponse()

			// Log based on error type
			if code >= 500 {
				logger.Error().
					Err(appErr.Err).
					Str("type", string(appErr.Type)).
					Str("message", appErr.Message).
					Str("path", c.Request().URL.Path).
					Str("method", c.Request().Method).
					Msg("Internal server error")
			} else if code >= 400 {
				logger.Warn().
					Str("type", string(appErr.Type)).
					Str("message", appErr.Message).
					Str("path", c.Request().URL.Path).
					Str("method", c.Request().Method).
					Msg("Client error")
			}
		} else if echoErr, ok := err.(*echo.HTTPError); ok {
			// Handle Echo's HTTP errors
			code = echoErr.Code
			message = map[string]interface{}{
				"type":    "HTTP_ERROR",
				"message": echoErr.Message,
			}

			if code >= 500 {
				logger.Error().
					Err(err).
					Int("code", code).
					Str("path", c.Request().URL.Path).
					Str("method", c.Request().Method).
					Msg("Echo HTTP error")
			}
		} else {
			// Unknown error
			code = http.StatusInternalServerError
			message = map[string]interface{}{
				"type":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			}

			logger.Error().
				Err(err).
				Str("path", c.Request().URL.Path).
				Str("method", c.Request().Method).
				Msg("Unexpected error")
		}

		// Send response
		if err := c.JSON(code, message); err != nil {
			logger.Error().Err(err).Msg("Failed to send error response")
		}
	}
}
