package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Harmeet10000/Fortress_API/src/internal/app"

	"github.com/labstack/echo/v4"
)

type OpenAPIHandler struct {
	Handler
}

func NewOpenAPIHandler(s *app.Server) *OpenAPIHandler {
	return &OpenAPIHandler{
		Handler: NewHandler(s),
	}
}

func (h *OpenAPIHandler) ServeOpenAPIUI(c echo.Context) error {
	templatePath := filepath.Join("src", "static", "openapi.html")
	templateBytes, err := os.ReadFile(templatePath)
	c.Response().Header().Set("Cache-Control", "no-cache")
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI UI template: %w", err)
	}

	templateString := string(templateBytes)

	err = c.HTML(http.StatusOK, templateString)
	if err != nil {
		return fmt.Errorf("failed to write HTML response: %w", err)
	}

	return nil
}
