package todo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *Repository) Create(ctx context.Context, req *CreateTodoRequest) (*Todo, error) {
	var description, categoryID pgtype.Text
	var dueDate pgtype.Timestamptz
	status := "pending"
	priority := "medium"

	if req.Description != nil {
		description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.CategoryID != nil {
		categoryID = pgtype.Text{String: req.CategoryID.String(), Valid: true}
	}
	if req.DueDate != nil {
		dueDate = pgtype.Timestamptz{Time: *req.DueDate, Valid: true}
	}
	if req.Status != nil {
		status = *req.Status
	}
	if req.Priority != nil {
		priority = *req.Priority
	}

	result, err := r.queries.CreateTodo(ctx, db.CreateTodoParams{
		Title:       req.Title,
		Description: description,
		Status:      db.TodoStatus(status),
		Priority:    db.TodoPriority(priority),
		CategoryID:  uuid.NullUUID{UUID: *req.CategoryID, Valid: req.CategoryID != nil},
		DueDate:     dueDate,
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to create todo")
		return nil, errs.NewInternalError("Failed to create todo", err)
	}

	return r.toModel(&result), nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Todo, error) {
	result, err := r.queries.GetTodoByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NewNotFoundError("Todo")
		}
		r.logger.Error().Err(err).Msg("Failed to get todo")
		return nil, errs.NewInternalError("Failed to get todo", err)
	}
	return r.toModelWithCategory(&result), nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Todo, error) {
	results, err := r.queries.ListTodos(ctx, db.ListTodosParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list todos")
		return nil, errs.NewInternalError("Failed to list todos", err)
	}

	todos := make([]Todo, 0, len(results))
	for _, result := range results {
		todos = append(todos, *r.toModelWithCategory(&result))
	}
	return todos, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountTodos(ctx)
	if err != nil {
		return 0, errs.NewInternalError("Failed to count todos", err)
	}
	return count, nil
}

func (r *Repository) toModel(dbTodo *db.Todo) *Todo {
	todo := &Todo{
		ID:        dbTodo.ID,
		Title:     dbTodo.Title,
		Status:    TodoStatus(dbTodo.Status),
		Priority:  TodoPriority(dbTodo.Priority),
		CreatedAt: dbTodo.CreatedAt.Time,
		UpdatedAt: dbTodo.UpdatedAt.Time,
	}
	if dbTodo.Description.Valid {
		todo.Description = &dbTodo.Description.String
	}
	if dbTodo.CategoryID.Valid {
		todo.CategoryID = &dbTodo.CategoryID.UUID
	}
	if dbTodo.DueDate.Valid {
		todo.DueDate = &dbTodo.DueDate.Time
	}
	if dbTodo.CompletedAt.Valid {
		todo.CompletedAt = &dbTodo.CompletedAt.Time
	}
	return todo
}

func (r *Repository) toModelWithCategory(row *db.GetTodoByIDRow) *Todo {
	todo := &Todo{
		ID:        row.ID,
		Title:     row.Title,
		Status:    TodoStatus(row.Status),
		Priority:  TodoPriority(row.Priority),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.Description.Valid {
		todo.Description = &row.Description.String
	}
	if row.CategoryID.Valid {
		todo.CategoryID = &row.CategoryID.UUID
	}
	if row.CategoryName.Valid {
		todo.CategoryName = &row.CategoryName.String
	}
	return todo
}
