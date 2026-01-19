package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/errs"
)

// Service handles category business logic
type Service struct {
	repo   *Repository
	logger *zerolog.Logger
}

// NewService creates a new category service
func NewService(repo *Repository, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new category
func (s *Service) Create(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if category with same name already exists
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil && !errs.IsAppError(err) {
		return nil, err
	}
	if existing != nil {
		return nil, errs.NewConflictError("Category with this name already exists")
	}

	// Create category
	category, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	s.logger.Info().
		Str("category_id", category.ID.String()).
		Str("name", category.Name).
		Msg("Category created successfully")

	return s.toResponse(category), nil
}

// GetByID retrieves a category by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(category), nil
}

// List retrieves a paginated list of categories
func (s *Service) List(ctx context.Context, req *ListCategoriesRequest) (*PaginatedCategoriesResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get categories
	categories, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to responses
	responses := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, *s.toResponse(&category))
	}

	return &PaginatedCategoriesResponse{
		Data: responses,
		Pagination: PaginationMeta{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}, nil
}

// Update updates a category
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if category exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check name uniqueness if updating name
	if req.Name != nil {
		existing, err := s.repo.GetByName(ctx, *req.Name)
		if err != nil && !errs.IsAppError(err) {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errs.NewConflictError("Category with this name already exists")
		}
	}

	// Update category
	category, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	s.logger.Info().
		Str("category_id", category.ID.String()).
		Msg("Category updated successfully")

	return s.toResponse(category), nil
}

// Delete deletes a category
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if category exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete category
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.Info().
		Str("category_id", id.String()).
		Msg("Category deleted successfully")

	return nil
}

// toResponse converts domain model to response DTO
func (s *Service) toResponse(category *Category) *CategoryResponse {
	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
