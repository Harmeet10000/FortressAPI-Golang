package category

import (
	"github.com/google/uuid"
	"github.com/yourusername/task-management-api/internal/validation"
)

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// Validate validates the create category request
func (r *CreateCategoryRequest) Validate() error {
	v := validation.NewValidator()
	
	v.Required("name", r.Name)
	v.MinLength("name", r.Name, 1)
	v.MaxLength("name", r.Name, 100)

	if r.Color != nil && *r.Color != "" {
		v.Custom("color", len(*r.Color) == 7 && (*r.Color)[0] == '#',
			"color must be a valid hex color code (e.g., #FF5733)")
	}

	return v.Validate()
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// Validate validates the update category request
func (r *UpdateCategoryRequest) Validate() error {
	v := validation.NewValidator()

	if r.Name != nil {
		v.MinLength("name", *r.Name, 1)
		v.MaxLength("name", *r.Name, 100)
	}

	if r.Color != nil && *r.Color != "" {
		v.Custom("color", len(*r.Color) == 7 && (*r.Color)[0] == '#',
			"color must be a valid hex color code (e.g., #FF5733)")
	}

	return v.Validate()
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Color       *string   `json:"color,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// ListCategoriesRequest represents the request to list categories
type ListCategoriesRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// Validate validates the list categories request
func (r *ListCategoriesRequest) Validate() error {
	v := validation.NewValidator()

	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Offset < 0 {
		r.Offset = 0
	}

	v.Max("limit", r.Limit, 100)
	v.Min("offset", r.Offset, 0)

	return v.Validate()
}

// PaginatedCategoriesResponse represents paginated categories
type PaginatedCategoriesResponse struct {
	Data       []CategoryResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
