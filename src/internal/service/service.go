package service

import (
	"fmt"

	"github.com/Harmeet10000/Fortress_API/src/internal/helper/aws"
	"github.com/Harmeet10000/Fortress_API/src/internal/helper/job"
	"github.com/Harmeet10000/Fortress_API/src/internal/repository"
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
)

type Services struct {
	Auth     *AuthService
	Job      *job.JobService
	Todo     *TodoService
	Comment  *CommentService
	Category *CategoryService
}

func NewServices(s *app.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	s.Job.SetAuthService(authService)

	awsClient, err := aws.NewAWS(s)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	return &Services{
		Job:      s.Job,
		Auth:     authService,
		Category: NewCategoryService(s, repos.Category),
		Comment:  NewCommentService(s, repos.Comment, repos.Todo),
		Todo:     NewTodoService(s, repos.Todo, repos.Category, awsClient),
	}, nil
}
