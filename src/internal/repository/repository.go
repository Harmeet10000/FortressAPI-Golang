package repository

import "github.com/Harmeet10000/Fortress_API/src/internal/app"


type Repositories struct {
	Todo     *TodoRepository
	// Comment  *CommentRepository
	// Category *CategoryRepository
}

func NewRepositories(s *app.Server) *Repositories {
	return &Repositories{
		Todo:     NewTodoRepository(s),
		// Comment:  NewCommentRepository(s),
		// Category: NewCategoryRepository(s),
	}
}
