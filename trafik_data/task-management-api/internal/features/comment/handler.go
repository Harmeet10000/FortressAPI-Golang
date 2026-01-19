package comment

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	appErrors "github.com/yourusername/task-management-api/internal/errors"
	"github.com/yourusername/task-management-api/internal/validation"
)

type Handler struct {
	service   *Service
	validator *validation.Validator
	logger    *zerolog.Logger
}

func NewHandler(service *Service, validator *validation.Validator, logger *zerolog.Logger) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	comments := g.Group("/comments")

	comments.POST("", h.Create)
	comments.GET("/:id", h.GetByID)
	comments.PUT("/:id", h.Update)
	comments.DELETE("/:id", h.Delete)

	// Comments by todo
	g.GET("/todos/:todo_id/comments", h.ListByTodo)
}

// Create godoc
// @Summary Create a new comment
// @Description Create a new comment for a todo
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body CreateCommentRequest true "Comment details"
// @Success 201 {object} CommentResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /comments [post]
func (h *Handler) Create(c echo.Context) error {
	var req CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return appErrors.NewBadRequestError("Invalid request body", err)
	}

	if err := h.validator.Validate(req); err != nil {
		return err
	}

	comment, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, comment)
}

// GetByID godoc
// @Summary Get a comment by ID
// @Description Get a comment by its ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID" format(uuid)
// @Success 200 {object} CommentResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /comments/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return appErrors.NewBadRequestError("Invalid comment ID", err)
	}

	comment, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, comment)
}

// ListByTodo godoc
// @Summary List comments for a todo
// @Description Get a paginated list of comments for a specific todo
// @Tags comments
// @Accept json
// @Produce json
// @Param todo_id path string true "Todo ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} CommentListResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /todos/{todo_id}/comments [get]
func (h *Handler) ListByTodo(c echo.Context) error {
	todoID, err := uuid.Parse(c.Param("todo_id"))
	if err != nil {
		return appErrors.NewBadRequestError("Invalid todo ID", err)
	}

	var page, pageSize int
	if err := echo.QueryParamsBinder(c).
		Int("page", &page).
		Int("page_size", &pageSize).
		BindError(); err != nil {
		return appErrors.NewBadRequestError("Invalid query parameters", err)
	}

	comments, err := h.service.ListByTodo(c.Request().Context(), todoID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, comments)
}

// Update godoc
// @Summary Update a comment
// @Description Update a comment by its ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID" format(uuid)
// @Param comment body UpdateCommentRequest true "Comment details"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /comments/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return appErrors.NewBadRequestError("Invalid comment ID", err)
	}

	var req UpdateCommentRequest
	if err := c.Bind(&req); err != nil {
		return appErrors.NewBadRequestError("Invalid request body", err)
	}

	if err := h.validator.Validate(req); err != nil {
		return err
	}

	comment, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, comment)
}

// Delete godoc
// @Summary Delete a comment
// @Description Delete a comment by its ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID" format(uuid)
// @Success 204
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /comments/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return appErrors.NewBadRequestError("Invalid comment ID", err)
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
