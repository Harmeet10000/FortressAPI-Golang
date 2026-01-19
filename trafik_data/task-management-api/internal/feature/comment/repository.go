package comment

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/connections"
	"github.com/yourusername/task-management-api/internal/database/db"
	"github.com/yourusername/task-management-api/internal/errs"
)

type Repository struct {
	db      *connections.Database
	queries *db.Queries
	logger  *zerolog.Logger
}

func NewRepository(database *connections.Database, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:      database,
		queries: db.New(database.Pool),
		logger:  logger,
	}
}

func (r *Repository) Create(ctx context.Context, req *CreateCommentRequest) (*Comment, error) {
	result, err := r.queries.CreateComment(ctx, db.CreateCommentParams{
		TodoID:  req.TodoID,
		Content: req.Content,
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to create comment")
		return nil, errs.NewInternalError("Failed to create comment", err)
	}
	return r.toModel(&result), nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Comment, error) {
	result, err := r.queries.GetCommentByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewNotFoundError("Comment")
		}
		r.logger.Error().Err(err).Msg("Failed to get comment")
		return nil, errs.NewInternalError("Failed to get comment", err)
	}
	return r.toModel(&result), nil
}

func (r *Repository) ListByTodoID(ctx context.Context, todoID uuid.UUID, limit, offset int) ([]Comment, error) {
	results, err := r.queries.ListCommentsByTodoID(ctx, db.ListCommentsByTodoIDParams{
		TodoID: todoID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list comments")
		return nil, errs.NewInternalError("Failed to list comments", err)
	}

	comments := make([]Comment, 0, len(results))
	for _, result := range results {
		comments = append(comments, *r.toModel(&result))
	}
	return comments, nil
}

func (r *Repository) CountByTodoID(ctx context.Context, todoID uuid.UUID) (int64, error) {
	count, err := r.queries.CountCommentsByTodoID(ctx, todoID)
	if err != nil {
		return 0, errs.NewInternalError("Failed to count comments", err)
	}
	return count, nil
}

func (r *Repository) toModel(dbComment *db.Comment) *Comment {
	return &Comment{
		ID:        dbComment.ID,
		TodoID:    dbComment.TodoID,
		Content:   dbComment.Content,
		CreatedAt: dbComment.CreatedAt.Time,
		UpdatedAt: dbComment.UpdatedAt.Time,
	}
}
