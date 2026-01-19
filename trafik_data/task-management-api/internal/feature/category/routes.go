package category

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers category routes
func RegisterRoutes(e *echo.Echo, handler *Handler) {
	categories := e.Group("/api/v1/categories")

	categories.POST("", handler.Create)
	categories.GET("", handler.List)
	categories.GET("/:id", handler.GetByID)
	categories.PUT("/:id", handler.Update)
	categories.DELETE("/:id", handler.Delete)
}
