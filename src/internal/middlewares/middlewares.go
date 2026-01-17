package middlewares

import (
	"github.com/Harmeet10000/Fortress_API/src/internal/app"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middlewares struct {
	Global          *GlobalMiddlewares
	// Auth            *AuthMiddleware
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	// RateLimit       *RateLimitMiddleware
}

func NewMiddlewares(s *app.Server) *Middlewares {
	// Get New Relic application instance from server
	var nrApp *newrelic.Application
	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}

	return &Middlewares{
		Global:          NewGlobalMiddlewares(s),
		// Auth:            NewAuthMiddleware(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracingMiddleware(s, nrApp),
		// RateLimit:       NewRateLimitMiddleware(s),
	}
}
