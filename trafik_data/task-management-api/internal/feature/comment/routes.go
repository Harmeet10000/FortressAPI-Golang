package comment

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	comments := e.Group("/api/v1/comments")
	comments.POST("", handler.Create)
	comments.GET("/:id", handler.GetByID)
	
	todos := e.Group("/api/v1/todos")
	todos.GET("/:todoId/comments", handler.ListByTodoID)
}
