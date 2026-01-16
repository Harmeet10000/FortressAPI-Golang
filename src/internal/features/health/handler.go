package health

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

)

// HealthController handles health check endpoints
type HealthController struct {
	db     *pgx.Pool
	redis  *redis.Client
	logger *zap.Logger
}

// NewHealthController creates a new health controller
func NewHealthController(db *pgx.Pool, redis *redis.Client, logger *zap.Logger) *HealthController {
	return &HealthController{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}


// Health checks the overall health of the application
// @Summary Get overall application health
// @Description Returns comprehensive health status of the application including system, database, and cache
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (hc *HealthController) Health(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hc.logger.Info("health check started")

	now := time.Now()

	// Get system health
	systemHealth := GetSystemHealth()

	// Get application health
	appHealth := GetApplicationHealth()

	// Check database
	dbHealth := CheckDatabasePool(ctx, hc.db)

	// Check Redis
	redisHealth := CheckRedis(ctx, hc.redis)

	// Check memory
	memHealth := CheckMemory()

	// Check disk
	diskHealth := CheckDisk()

	// Get CPU info
	cpuInfo := CheckCPU()

	// Determine overall status
	overallStatus := "healthy"
	checks := map[string]string{
		"database": dbHealth.Status,
		"redis":    redisHealth.Status,
		"memory":   memHealth.Status,
		"disk":     diskHealth.Status,
	}

	for _, check := range checks {
		if check == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
		if check == "warning" && overallStatus != "unhealthy" {
			overallStatus = "warning"
		}
	}

	response := HealthResponse{
		Status:        overallStatus,
		Timestamp:     now.Format(time.RFC3339),
		System:        systemHealth,
		Application:   appHealth,
		Database:      dbHealth,
		Redis:         redisHealth,
		Memory:        memHealth,
		Disk:          diskHealth,
		CPU:           cpuInfo,
		Checks:        checks,
	}

	hc.logger.Info("health check completed", zap.String("status", overallStatus))

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, response)
}

// LivenessProbe checks if the application is alive (minimal check)
// @Summary Liveness probe
// @Description Returns 200 if the application is running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (hc *HealthController) LivenessProbe(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// ReadinessProbe checks if the application is ready to handle requests
// @Summary Readiness probe
// @Description Returns 200 if the application is ready, 503 if not
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func (hc *HealthController) ReadinessProbe(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check database
	dbHealth := CheckDatabasePool(ctx, hc.db)
	if dbHealth.Status != "healthy" {
		hc.logger.Warn("readiness probe failed", zap.String("reason", "database"))
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":   "not_ready",
			"reason":   "database_unhealthy",
			"time":     time.Now().Format(time.RFC3339),
		})
	}

	// Check Redis
	redisHealth := CheckRedis(ctx, hc.redis)
	if redisHealth.Status != "healthy" {
		hc.logger.Warn("readiness probe failed", zap.String("reason", "redis"))
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":   "not_ready",
			"reason":   "redis_unhealthy",
			"time":     time.Now().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// SystemHealth returns system-level metrics
// @Summary Get system health metrics
// @Description Returns CPU, memory, and platform information
// @Tags Health
// @Produce json
// @Success 200 {object} SystemHealthResponse
// @Router /health/system [get]
func (hc *HealthController) SystemHealth(c echo.Context) error {
	system := GetSystemHealth()
	return c.JSON(http.StatusOK, system)
}

// ApplicationHealth returns application-level metrics
// @Summary Get application health metrics
// @Description Returns uptime, memory usage, and version information
// @Tags Health
// @Produce json
// @Success 200 {object} ApplicationHealthResponse
// @Router /health/app [get]
func (hc *HealthController) ApplicationHealth(c echo.Context) error {
	app := GetApplicationHealth()
	return c.JSON(http.StatusOK, app)
}

// MemoryHealth returns memory usage information
// @Summary Get memory health status
// @Description Returns memory usage percentage and status
// @Tags Health
// @Produce json
// @Success 200 {object} MemoryHealthResponse
// @Router /health/memory [get]
func (hc *HealthController) MemoryHealth(c echo.Context) error {
	memory := CheckMemory()
	return c.JSON(http.StatusOK, memory)
}
