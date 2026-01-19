package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/database"
	appErrors "github.com/yourusername/task-management-api/internal/errors"
	"github.com/yourusername/task-management-api/internal/features/todo"
)

type Service struct {
	repo     *Repository
	todoRepo *todo.Repository
	logger   *zerolog.Logger
}

func NewService(repo *Repository, todoRepo *todo.Repository, logger *zerolog.Logger) *Service {
	return &Service{
		repo:     repo,
		todoRepo: todoRepo,
		logger:   logger,
	}
}

func (s *Service) Create(ctx context.Context, req CreateCommentRequest) (*CommentResponse, error) {
	// Verify todo exists
	_, err := s.todoRepo.GetByID(ctx, req.TodoID)
	if err != nil {
		return nil, err
	}

	comment, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	response := ToCommentResponse(comment)
	return &response, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*CommentResponse, error) {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := ToCommentResponse(comment)
	return &response, nil
}

func (s *Service) ListByTodo(ctx context.Context, todoID uuid.UUID, page, pageSize int) (*CommentListResponse, error) {
	// Verify todo exists
	_, err := s.todoRepo.GetByID(ctx, todoID)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	comments, total, err := s.repo.ListByTodo(ctx, todoID, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]CommentResponse, len(comments))
	for i, c := range comments {
		responses[i] = ToCommentResponse(&c)
	}

	return &CommentListResponse{
		Comments: responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateCommentRequest) (*CommentResponse, error) {
	// Verify comment exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comment, err := s.repo.Update(ctx, id, req.Content)
	if err != nil {
		return nil, err
	}

	response := ToCommentResponse(comment)
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Verify comment exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func ToCommentResponse(comment *database.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		TodoID:    comment.TodoID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Time,
		UpdatedAt: comment.UpdatedAt.Time,
	}
}
