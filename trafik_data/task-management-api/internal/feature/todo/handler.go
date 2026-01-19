package todo

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/errs"
)

type Handler struct {
	service *Service
	logger  *zerolog.Logger
}

func NewHandler(service *Service, logger *zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid request body")
	}
	todo, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, todo)
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid todo ID")
	}
	todo, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func (h *Handler) List(c echo.Context) error {
	var req ListTodosRequest
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid query parameters")
	}
	todos, err := h.service.List(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todos)
}
