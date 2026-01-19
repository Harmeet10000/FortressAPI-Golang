package handler

import (
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
	"github.com/Harmeet10000/Fortress_API/src/internal/service"
)

type Handlers struct {
	Health   *HealthHandler
	OpenAPI  *OpenAPIHandler
	Todo     *TodoHandler
	Comment  *CommentHandler
	Category *CategoryHandler
}

func NewHandlers(s *app.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(s),
		OpenAPI:  NewOpenAPIHandler(s),
		Todo:     NewTodoHandler(s, services.Todo),
		Category: NewCategoryHandler(s, services.Category),
		Comment:  NewCommentHandler(s, services.Comment),
	}
}
