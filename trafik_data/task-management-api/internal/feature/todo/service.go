package todo

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/feature/category"
)

type Service struct {
	repo         *Repository
	categoryRepo *category.Repository
	logger       *zerolog.Logger
}

func NewService(repo *Repository, categoryRepo *category.Repository, logger *zerolog.Logger) *Service {
	return &Service{
		repo:         repo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateTodoRequest) (*TodoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.CategoryID != nil {
		if _, err := s.categoryRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, err
		}
	}

	todo, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	s.logger.Info().Str("todo_id", todo.ID.String()).Msg("Todo created")
	return s.toResponse(todo), nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*TodoResponse, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(todo), nil
}

func (s *Service) List(ctx context.Context, req *ListTodosRequest) (*PaginatedTodosResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	todos, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]TodoResponse, 0, len(todos))
	for _, todo := range todos {
		responses = append(responses, *s.toResponse(&todo))
	}

	return &PaginatedTodosResponse{
		Data: responses,
		Pagination: PaginationMeta{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}, nil
}

func (s *Service) toResponse(todo *Todo) *TodoResponse {
	resp := &TodoResponse{
		ID:           todo.ID,
		Title:        todo.Title,
		Description:  todo.Description,
		Status:       string(todo.Status),
		Priority:     string(todo.Priority),
		CategoryID:   todo.CategoryID,
		CategoryName: todo.CategoryName,
		CreatedAt:    todo.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    todo.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if todo.DueDate != nil {
		dueDate := todo.DueDate.Format("2006-01-02T15:04:05Z07:00")
		resp.DueDate = &dueDate
	}
	if todo.CompletedAt != nil {
		completedAt := todo.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.CompletedAt = &completedAt
	}
	return resp
}
