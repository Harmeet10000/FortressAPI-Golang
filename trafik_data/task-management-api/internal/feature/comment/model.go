package comment

import (
	"time"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id"`
	TodoID    uuid.UUID `json:"todo_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
