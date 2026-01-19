package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/task-management-api/internal/validation"
)

type CreateTodoRequest struct {
	Title       string        `json:"title"`
	Description *string       `json:"description,omitempty"`
	Status      *string       `json:"status,omitempty"`
	Priority    *string       `json:"priority,omitempty"`
	CategoryID  *uuid.UUID    `json:"category_id,omitempty"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
}

func (r *CreateTodoRequest) Validate() error {
	v := validation.NewValidator()
	v.Required("title", r.Title)
	v.MinLength("title", r.Title, 1)
	v.MaxLength("title", r.Title, 255)

	if r.Status != nil {
		v.In("status", *r.Status, []string{"pending", "in_progress", "completed"})
	}
	if r.Priority != nil {
		v.In("priority", *r.Priority, []string{"low", "medium", "high"})
	}

	return v.Validate()
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

func (r *UpdateTodoRequest) Validate() error {
	v := validation.NewValidator()
	if r.Title != nil {
		v.MinLength("title", *r.Title, 1)
		v.MaxLength("title", *r.Title, 255)
	}
	if r.Status != nil {
		v.In("status", *r.Status, []string{"pending", "in_progress", "completed"})
	}
	if r.Priority != nil {
		v.In("priority", *r.Priority, []string{"low", "medium", "high"})
	}
	return v.Validate()
}

type UpdateTodoStatusRequest struct {
	Status string `json:"status"`
}

func (r *UpdateTodoStatusRequest) Validate() error {
	v := validation.NewValidator()
	v.Required("status", r.Status)
	v.In("status", r.Status, []string{"pending", "in_progress", "completed"})
	return v.Validate()
}

type TodoResponse struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	CategoryName *string    `json:"category_name,omitempty"`
	DueDate      *string    `json:"due_date,omitempty"`
	CompletedAt  *string    `json:"completed_at,omitempty"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

type ListTodosRequest struct {
	Status     *string `query:"status"`
	CategoryID *string `query:"category_id"`
	Limit      int     `query:"limit"`
	Offset     int     `query:"offset"`
}

func (r *ListTodosRequest) Validate() error {
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

type PaginatedTodosResponse struct {
	Data       []TodoResponse `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
