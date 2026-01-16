package health

// HealthResponse represents the overall health status response
type HealthResponse struct {
	Status              string                                `json:"status"`
	Timestamp           string                                `json:"timestamp"`
	System              SystemHealthResponse            `json:"system"`
	Application         ApplicationHealthResponse       `json:"application"`
	Database            HealthCheckResponse             `json:"database"`
	Redis               HealthCheckResponse             `json:"redis"`
	Memory              MemoryHealthResponse            `json:"memory"`
	Disk                DiskHealthResponse              `json:"disk"`
	CPU                 map[string]interface{}                `json:"cpu"`
	Checks              map[string]string                     `json:"checks"`
}

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
