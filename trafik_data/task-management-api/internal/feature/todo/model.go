package todo

import (
	"time"

	"github.com/google/uuid"
)

// TodoStatus represents the status of a todo
type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted  TodoStatus = "completed"
)

// TodoPriority represents the priority of a todo
type TodoPriority string

const (
	TodoPriorityLow    TodoPriority = "low"
	TodoPriorityMedium TodoPriority = "medium"
	TodoPriorityHigh   TodoPriority = "high"
)

// Todo represents a task
type Todo struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title"`
	Description  *string       `json:"description,omitempty"`
	Status       TodoStatus    `json:"status"`
	Priority     TodoPriority  `json:"priority"`
	CategoryID   *uuid.UUID    `json:"category_id,omitempty"`
	CategoryName *string       `json:"category_name,omitempty"`
	DueDate      *time.Time    `json:"due_date,omitempty"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
