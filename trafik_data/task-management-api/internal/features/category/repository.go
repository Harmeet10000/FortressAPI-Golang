package category

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/connections"
	"github.com/yourusername/task-management-api/internal/database/db"
	"github.com/yourusername/task-management-api/internal/errs"
)

// Repository handles category data persistence
type Repository struct {
	db      *connections.Database
	queries *db.Queries
	logger  *zerolog.Logger
}

// NewRepository creates a new category repository
func NewRepository(database *connections.Database, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:      database,
		queries: db.New(database.Pool),
		logger:  logger,
	}
}

// Create creates a new category
func (r *Repository) Create(ctx context.Context, req *CreateCategoryRequest) (*Category, error) {
	var description, color pgtype.Text
	if req.Description != nil {
		description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.Color != nil {
		color = pgtype.Text{String: *req.Color, Valid: true}
	}

	result, err := r.queries.CreateCategory(ctx, db.CreateCategoryParams{
		Name:        req.Name,
		Description: description,
		Color:       color,
	})
	if err != nil {
		r.logger.Error().Err(err).Str("name", req.Name).Msg("Failed to create category")
		return nil, errs.NewInternalError("Failed to create category", err)
	}

	return r.toModel(&result), nil
}

// GetByID retrieves a category by ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	result, err := r.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewNotFoundError("Category")
		}
		r.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to get category")
		return nil, errs.NewInternalError("Failed to get category", err)
	}

	return r.toModel(&result), nil
}

// GetByName retrieves a category by name
func (r *Repository) GetByName(ctx context.Context, name string) (*Category, error) {
	result, err := r.queries.GetCategoryByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewNotFoundError("Category")
		}
		r.logger.Error().Err(err).Str("name", name).Msg("Failed to get category by name")
		return nil, errs.NewInternalError("Failed to get category", err)
	}

	return r.toModel(&result), nil
}

// List retrieves a paginated list of categories
func (r *Repository) List(ctx context.Context, limit, offset int) ([]Category, error) {
	results, err := r.queries.ListCategories(ctx, db.ListCategoriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list categories")
		return nil, errs.NewInternalError("Failed to list categories", err)
	}

	categories := make([]Category, 0, len(results))
	for _, result := range results {
		categories = append(categories, *r.toModel(&result))
	}

	return categories, nil
}

// Update updates a category
func (r *Repository) Update(ctx context.Context, id uuid.UUID, req *UpdateCategoryRequest) (*Category, error) {
	var name, description, color pgtype.Text

	if req.Name != nil {
		name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.Description != nil {
		description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.Color != nil {
		color = pgtype.Text{String: *req.Color, Valid: true}
	}

	result, err := r.queries.UpdateCategory(ctx, db.UpdateCategoryParams{
		ID:          id,
		Name:        name,
		Description: description,
		Color:       color,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewNotFoundError("Category")
		}
		r.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to update category")
		return nil, errs.NewInternalError("Failed to update category", err)
	}

	return r.toModel(&result), nil
}

// Delete deletes a category
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteCategory(ctx, id)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to delete category")
		return errs.NewInternalError("Failed to delete category", err)
	}

	return nil
}

// Count counts total categories
func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountCategories(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to count categories")
		return 0, errs.NewInternalError("Failed to count categories", err)
	}

	return count, nil
}

// toModel converts database model to domain model
func (r *Repository) toModel(dbCategory *db.Category) *Category {
	category := &Category{
		ID:        dbCategory.ID,
		Name:      dbCategory.Name,
		CreatedAt: dbCategory.CreatedAt.Time,
		UpdatedAt: dbCategory.UpdatedAt.Time,
	}

	if dbCategory.Description.Valid {
		category.Description = &dbCategory.Description.String
	}
	if dbCategory.Color.Valid {
		category.Color = &dbCategory.Color.String
	}

	return category
}
