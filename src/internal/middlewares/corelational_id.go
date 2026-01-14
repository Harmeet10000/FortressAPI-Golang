package middlewares

import (
	"context"
	"net/http"

	"github.com/Harmeet10000/Fortress_API/src/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CorrelationIDMiddleware adds a correlation ID to the request context
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}

		reqLogger := utils.GetLogger().With(
			zap.String("correlation_id", correlationID),
		)

		ctx := context.WithValue(r.Context(), utils.RequestLoggerKey, reqLogger)
		ctx = context.WithValue(ctx, utils.RequestIDKey, correlationID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// generateCorrelationID generates a unique correlation ID
func generateCorrelationID() string {
	return uuid.New().String()
}
