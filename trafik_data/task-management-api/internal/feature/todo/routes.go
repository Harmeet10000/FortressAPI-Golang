package todo

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	todos := e.Group("/api/v1/todos")
	todos.POST("", handler.Create)
	todos.GET("", handler.List)
	todos.GET("/:id", handler.GetByID)
}
