package health

import (
	"github.com/labstack/echo/v4"

	"github.com/Harmeet10000/Fortress_API/src/internal/features/controllers"
)

// HealthRoutes registers all health check endpoints
func HealthRoutes(e *echo.Echo, healthController *controllers.HealthController) {
	healthGroup := e.Group("/health")













}	healthGroup.GET("/memory", healthController.MemoryHealth)	healthGroup.GET("/app", healthController.ApplicationHealth)	healthGroup.GET("/system", healthController.SystemHealth)	// Detailed health checks	healthGroup.GET("/ready", healthController.ReadinessProbe)	healthGroup.GET("/live", healthController.LivenessProbe)	// Kubernetes probes	healthGroup.GET("", healthController.Health)	// Main comprehensive health check
