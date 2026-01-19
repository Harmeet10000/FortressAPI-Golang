package comment

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
	var req CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid request body")
	}
	comment, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid comment ID")
	}
	comment, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) ListByTodoID(c echo.Context) error {
	todoID, err := uuid.Parse(c.Param("todoId"))
	if err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid todo ID")
	}
	
	var req ListCommentsRequest
	req.TodoID = todoID
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid query parameters")
	}
	
	comments, err := h.service.ListByTodoID(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, comments)
}
