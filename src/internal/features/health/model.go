package health

import (
		"github.com/Harmeet10000/Fortress_API/src/internal/utils"
)

// HealthResponse represents the overall health status response
type HealthResponse struct {
	Status              string                                `json:"status"`
	Timestamp           string                                `json:"timestamp"`
	System              utils.SystemHealthResponse            `json:"system"`
	Application         utils.ApplicationHealthResponse       `json:"application"`
	Database            utils.HealthCheckResponse             `json:"database"`
	Redis               utils.HealthCheckResponse             `json:"redis"`
	Memory              utils.MemoryHealthResponse            `json:"memory"`
	Disk                utils.DiskHealthResponse              `json:"disk"`
	CPU                 map[string]interface{}                `json:"cpu"`
	Checks              map[string]string                     `json:"checks"`
}
