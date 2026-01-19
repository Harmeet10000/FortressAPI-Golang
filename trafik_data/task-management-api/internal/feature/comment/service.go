package comment

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/yourusername/task-management-api/internal/feature/todo"
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

func (s *Service) Create(ctx context.Context, req *CreateCommentRequest) (*CommentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.todoRepo.GetByID(ctx, req.TodoID); err != nil {
		return nil, err
	}

	comment, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	s.logger.Info().Str("comment_id", comment.ID.String()).Msg("Comment created")
	return s.toResponse(comment), nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*CommentResponse, error) {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(comment), nil
}

func (s *Service) ListByTodoID(ctx context.Context, req *ListCommentsRequest) (*PaginatedCommentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	comments, err := s.repo.ListByTodoID(ctx, req.TodoID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByTodoID(ctx, req.TodoID)
	if err != nil {
		return nil, err
	}

	responses := make([]CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, *s.toResponse(&comment))
	}

	return &PaginatedCommentsResponse{
		Data: responses,
		Pagination: PaginationMeta{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}, nil
}

func (s *Service) toResponse(comment *Comment) *CommentResponse {
	return &CommentResponse{
		ID:        comment.ID,
		TodoID:    comment.TodoID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
