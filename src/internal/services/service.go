package services

import (
	// "github.com/Harmeet10000/Fortress_API/internal/lib/job"
	"github.com/Harmeet10000/Fortress_API/src/internal/repository"
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
)

type Services struct {
	Auth *AuthService
	// Job  *job.JobService
}

func NewServices(s *app.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		// Job:  s.Job,
		Auth: authService,
	}, nil
}
