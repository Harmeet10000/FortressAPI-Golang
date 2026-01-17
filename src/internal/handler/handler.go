package handler

import (
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
	"github.com/Harmeet10000/Fortress_API/src/internal/services"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *app.Server, services *services.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}
