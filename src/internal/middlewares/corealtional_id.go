package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

func CorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(RequestIDHeader)
			if correlationID == "" {
				correlationID = uuid.New().String() // 4c90fc3f-39cc-4b04-af21-c83ee64aa67e
			}

			c.Set(RequestIDKey, correlationID)
			c.Response().Header().Set(RequestIDHeader, correlationID)

			return next(c)
		}
	}
}

func GetCorrelationID(c echo.Context) string {
	if correlationID, ok := c.Get(RequestIDKey).(string); ok {
		return correlationID
	}
	return ""
}
