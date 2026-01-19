package comment

import (
	"github.com/google/uuid"
	"github.com/yourusername/task-management-api/internal/validation"
)

type CreateCommentRequest struct {
	TodoID  uuid.UUID `json:"todo_id"`
	Content string    `json:"content"`
}

func (r *CreateCommentRequest) Validate() error {
	v := validation.NewValidator()
	v.Required("content", r.Content)
	v.MinLength("content", r.Content, 1)
	v.MaxLength("content", r.Content, 1000)
	return v.Validate()
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

func (r *UpdateCommentRequest) Validate() error {
	v := validation.NewValidator()
	v.Required("content", r.Content)
	v.MinLength("content", r.Content, 1)
	v.MaxLength("content", r.Content, 1000)
	return v.Validate()
}

type CommentResponse struct {
	ID        uuid.UUID `json:"id"`
	TodoID    uuid.UUID `json:"todo_id"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type ListCommentsRequest struct {
	TodoID uuid.UUID `query:"todo_id"`
	Limit  int       `query:"limit"`
	Offset int       `query:"offset"`
}

func (r *ListCommentsRequest) Validate() error {
	v := validation.NewValidator()
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Offset < 0 {
		r.Offset = 0
	}
	v.Max("limit", r.Limit, 100)
	return v.Validate()
}

type PaginatedCommentsResponse struct {
	Data       []CommentResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}

type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
