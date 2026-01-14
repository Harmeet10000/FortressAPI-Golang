package health

import (
	"github.com/labstack/echo/v4"
)

// HealthRoutes registers all health check endpoints
func HealthRoutes(e *echo.Echo, healthController *HealthController) {
	healthGroup := e.Group("/health")

	// Main comprehensive health check
	healthGroup.GET("", healthController.Health)

	// Kubernetes probes
	healthGroup.GET("/live", healthController.LivenessProbe)
	healthGroup.GET("/ready", healthController.ReadinessProbe)

	// Detailed health checks
	healthGroup.GET("/system", healthController.SystemHealth)
	healthGroup.GET("/app", healthController.ApplicationHealth)
	healthGroup.GET("/memory", healthController.MemoryHealth)
}



