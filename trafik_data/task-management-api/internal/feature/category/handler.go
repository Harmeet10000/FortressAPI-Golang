package category

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/errs"
)

// Handler handles HTTP requests for categories
type Handler struct {
	service *Service
	logger  *zerolog.Logger
}

// NewHandler creates a new category handler
func NewHandler(service *Service, logger *zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create handles POST /categories
// @Summary Create a new category
// @Description Create a new category for organizing tasks
// @Tags categories
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Category creation request"
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 409 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /categories [post]
func (h *Handler) Create(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid request body")
	}

	category, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

// GetByID handles GET /categories/:id
// @Summary Get a category by ID
// @Description Retrieve a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /categories/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid category ID")
	}

	category, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

// List handles GET /categories
// @Summary List categories
// @Description Retrieve a paginated list of categories
// @Tags categories
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} PaginatedCategoriesResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /categories [get]
func (h *Handler) List(c echo.Context) error {
	var req ListCategoriesRequest
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid query parameters")
	}

	categories, err := h.service.List(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

// Update handles PUT /categories/:id
// @Summary Update a category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body UpdateCategoryRequest true "Category update request"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 409 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /categories/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid category ID")
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid request body")
	}

	category, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

// Delete handles DELETE /categories/:id
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /categories/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.New(errs.ErrorTypeBadRequest, "Invalid category ID")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
