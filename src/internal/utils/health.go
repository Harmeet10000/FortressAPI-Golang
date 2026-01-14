package utils

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// SystemHealthResponse represents system-level health metrics
type SystemHealthResponse struct {
	CPUUsage        []float64 `json:"cpu_usage"`
	CPUUsagePercent string    `json:"cpu_usage_percent"`
	TotalMemory     string    `json:"total_memory"`
	FreeMemory      string    `json:"free_memory"`
	Platform        string    `json:"platform"`
	Arch            string    `json:"arch"`
}

// ApplicationHealthResponse represents application-level health metrics
type ApplicationHealthResponse struct {
	Environment string                 `json:"environment"`
	Uptime      string                 `json:"uptime"`
	MemoryUsage ApplicationMemoryUsage `json:"memory_usage"`
	PID         int                    `json:"pid"`
	GoVersion   string                 `json:"go_version"`
}

// ApplicationMemoryUsage represents Go memory statistics
type ApplicationMemoryUsage struct {
	AllocMB      string `json:"alloc_mb"`
	TotalAllocMB string `json:"total_alloc_mb"`
	SysMB        string `json:"sys_mb"`
	NumGC        uint32 `json:"num_gc"`
	HeapAllocMB  string `json:"heap_alloc_mb"`
	HeapSysMB    string `json:"heap_sys_mb"`
	StackAllocMB string `json:"stack_alloc_mb"`
}

// HealthCheckResponse represents the result of a health check
type HealthCheckResponse struct {
	Status       string                 `json:"status"`
	ResponseTime int64                  `json:"response_time_ms"`
	Error        string                 `json:"error,omitempty"`
	Details      map[string]interface{} `json:"details,omitempty"`
}

// MemoryHealthResponse represents memory health status
type MemoryHealthResponse struct {
	Status       string `json:"status"`
	TotalMB      int    `json:"total_mb"`
	UsedMB       int    `json:"used_mb"`
	UsagePercent int    `json:"usage_percent"`
}

// DiskHealthResponse represents disk health status
type DiskHealthResponse struct {
	Status     string `json:"status"`
	Accessible bool   `json:"accessible"`
	Error      string `json:"error,omitempty"`
}

// startTime tracks application startup time
var startTime = time.Now()

// GetSystemHealth returns system-level health metrics
func GetSystemHealth() SystemHealthResponse {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get CPU load average (Unix-like systems)
	var cpuUsage [3]float64
	runtime.NumGoroutine() // Placeholder for cross-platform support

	// For cross-platform support, we use basic metrics
	// On Linux, load average is available via /proc/loadavg
	cpuLoadAvg := make([]float64, 0)

	// Calculate CPU usage percentage (approximate)
	numCPU := float64(runtime.NumCPU())
	cpuPercent := "N/A"

	totalMemMB := float64(m.Sys) / 1024 / 1024
	freeMemMB := float64(m.Sys-m.Alloc) / 1024 / 1024

	return SystemHealthResponse{
		CPUUsage:        cpuLoadAvg,
		CPUUsagePercent: cpuPercent,
		TotalMemory:     fmt.Sprintf("%.2f MB", totalMemMB),
		FreeMemory:      fmt.Sprintf("%.2f MB", freeMemMB),
		Platform:        os.Getenv("GOOS"),
		Arch:            runtime.GOARCH,
	}
}

// GetApplicationHealth returns application-level health metrics
func GetApplicationHealth() ApplicationHealthResponse {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime)
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	return ApplicationHealthResponse{
		Environment: env,
		Uptime:      fmt.Sprintf("%.2f seconds", uptime.Seconds()),
		MemoryUsage: ApplicationMemoryUsage{
			AllocMB:      fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
			TotalAllocMB: fmt.Sprintf("%.2f MB", float64(m.TotalAlloc)/1024/1024),
			SysMB:        fmt.Sprintf("%.2f MB", float64(m.Sys)/1024/1024),
			NumGC:        m.NumGC,
			HeapAllocMB:  fmt.Sprintf("%.2f MB", float64(m.HeapAlloc)/1024/1024),
			HeapSysMB:    fmt.Sprintf("%.2f MB", float64(m.HeapSys)/1024/1024),
			StackAllocMB: fmt.Sprintf("%.2f MB", float64(m.StackInuse)/1024/1024),
		},
		PID:       os.Getpid(),
		GoVersion: runtime.Version(),
	}
}

// CheckDatabase checks PostgreSQL connection health
func CheckDatabase(ctx context.Context, db *pgx.Conn) HealthCheckResponse {
	start := time.Now()
	response := HealthCheckResponse{
		Details: make(map[string]interface{}),
	}

	// Create a context with timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Test the connection with a simple ping
	err := db.Ping(pingCtx)
	responseTime := time.Since(start).Milliseconds()
	response.ResponseTime = responseTime

	if err != nil {
		response.Status = "unhealthy"
		response.Error = err.Error()
		response.Details["connection"] = "failed"
		return response
	}

	response.Status = "healthy"
	response.Details["connection"] = "connected"
	response.Details["driver"] = "pgx"
	response.Details["response_time_ms"] = responseTime

	return response
}

// CheckDatabasePool checks PostgreSQL connection pool health
func CheckDatabasePool(ctx context.Context, pool *pgxpool.Pool) HealthCheckResponse {
	start := time.Now()
	response := HealthCheckResponse{
		Details: make(map[string]interface{}),
	}

	// Create a context with timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Test the connection with a simple ping
	err := pool.Ping(pingCtx)
	responseTime := time.Since(start).Milliseconds()
	response.ResponseTime = responseTime

	if err != nil {
		response.Status = "unhealthy"
		response.Error = err.Error()
		response.Details["connection"] = "failed"
		return response
	}

	// Get pool stats
	stats := pool.Stat()
	response.Status = "healthy"
	response.Details["connection"] = "connected"
	response.Details["driver"] = "pgx"
	response.Details["response_time_ms"] = responseTime
	response.Details["acquired_conns"] = stats.AcquiredConns()
	response.Details["idle_conns"] = stats.IdleConns()
	response.Details["constructed_conns"] = stats.ConstructedConns()
	response.Details["max_conns"] = stats.MaxConns()

	return response
}

// CheckRedis checks Redis connection health
func CheckRedis(ctx context.Context, client *redis.Client) HealthCheckResponse {
	start := time.Now()
	response := HealthCheckResponse{
		Details: make(map[string]interface{}),
	}

	// Create a context with timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Ping Redis
	status := client.Ping(pingCtx)
	responseTime := time.Since(start).Milliseconds()
	response.ResponseTime = responseTime

	if status.Err() != nil {
		response.Status = "unhealthy"
		response.Error = status.Err().Error()
		response.Details["connection"] = "failed"
		return response
	}

	response.Status = "healthy"
	response.Details["connection"] = "connected"
	response.Details["message"] = status.Val()
	response.Details["response_time_ms"] = responseTime

	return response
}

// CheckMemory checks Go runtime memory usage and health
func CheckMemory() MemoryHealthResponse {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	totalMB := int(m.HeapAlloc / 1024 / 1024)
	usedMB := int(m.HeapAlloc / 1024 / 1024)
	usagePercent := 0
	if m.HeapAlloc > 0 {
		usagePercent = int((float64(m.HeapAlloc) / float64(m.HeapSys)) * 100)
	}

	status := "healthy"
	if usagePercent > 90 {
		status = "warning"
	} else if usagePercent > 95 {
		status = "critical"
	}

	return MemoryHealthResponse{
		Status:       status,
		TotalMB:      totalMB,
		UsedMB:       usedMB,
		UsagePercent: usagePercent,
	}
}

// CheckDisk checks if the filesystem is accessible
func CheckDisk() DiskHealthResponse {
	response := DiskHealthResponse{
		Accessible: false,
	}

	// Try to read the current working directory
	_, err := os.Stat(".")
	if err != nil {
		response.Status = "unhealthy"
		response.Error = err.Error()
		return response
	}

	// Try to check disk permissions by getting file info
	dir, err := os.Open(".")
	if err != nil {
		response.Status = "unhealthy"
		response.Error = err.Error()
		return response
	}
	defer dir.Close()

	response.Status = "healthy"
	response.Accessible = true
	return response
}

// CheckCPU returns basic CPU information (simplified)
func CheckCPU() map[string]interface{} {
	return map[string]interface{}{
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"arch":          runtime.GOARCH,
		"os":            runtime.GOOS,
		"compiler":      runtime.Compiler,
		"go_version":    runtime.Version(),
	}
}

// ResetUptime resets the application uptime counter (useful for testing)
func ResetUptime() {
	startTime = time.Now()
}
