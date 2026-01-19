package middleware

import (
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
)

type RateLimitMiddleware struct {
	server *app.Server
}

func NewRateLimitMiddleware(s *app.Server) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		server: s,
	}
}

func (r *RateLimitMiddleware) RecordRateLimitHit(endpoint string) {
	if r.server.LoggerService != nil && r.server.LoggerService.GetApplication() != nil {
		r.server.LoggerService.GetApplication().RecordCustomEvent("RateLimitHit", map[string]interface{}{
			"endpoint": endpoint,
		})
	}
}

