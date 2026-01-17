package services

import (
	"github.com/Harmeet10000/Fortress_API/src/internal/app"

	"github.com/clerk/clerk-sdk-go/v2"
)

type AuthService struct {
	server *app.Server
}

func NewAuthService(s *app.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server: s,
	}
}
