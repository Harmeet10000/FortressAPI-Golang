package router

import (
	"net/http"

	"github.com/Harmeet10000/Fortress_API/src/internal/app"
	"github.com/Harmeet10000/Fortress_API/src/internal/handler"
	"github.com/Harmeet10000/Fortress_API/src/internal/middlewares"
	"github.com/Harmeet10000/Fortress_API/src/internal/services"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"golang.org/x/time/rate"
)

func NewRouter(s *app.Server, h *handler.Handlers, services *services.Services) *echo.Echo {
	middlewares := middlewares.NewMiddlewares(s)

	router := echo.New()

	router.HTTPErrorHandler = middlewares.Global.GlobalErrorHandler

	// global middlewares
	router.Use(
		echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
			Store: echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20)),
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				// Record rate limit hit metrics
				if rateLimitMiddleware := middlewares.RateLimit; rateLimitMiddleware != nil {
					rateLimitMiddleware.RecordRateLimitHit(c.Path())
				}

				s.Logger.Warn().
					// Str("request_id", middlewares.GetCorrelationID(c)).
					Str("identifier", identifier).
					Str("path", c.Path()).
					Str("method", c.Request().Method).
					Str("ip", c.RealIP()).
					Msg("rate limit exceeded")

				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			},
		}),
		middlewares.Global.CORS(),
		middlewares.Global.Secure(),
		// middlewares.CorrelationID(),
		middlewares.Tracing.NewRelicMiddleware(),
		middlewares.Tracing.EnhanceTracing(),
		middlewares.ContextEnhancer.EnhanceContext(),
		middlewares.Global.RequestLogger(),
		middlewares.Global.Recover(),
	)

	// register system routes
	registerSystemRoutes(router, h)

	// register versioned routes
	router.Group("/api/v1")

	return router
}
