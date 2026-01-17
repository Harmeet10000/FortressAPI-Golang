package utils

import (
	"net/http"
	"net"
	"os"
	"strings"

	// "github.com/google/uuid" // optional — if you generate correlation IDs
)

// EnvProduction is the value that signals production mode
const EnvProduction = "production"

// APIResponse is the standard shape for all JSON API responses
type APIResponse[T any] struct {
	Success    bool              `json:"success"`
	StatusCode int               `json:"statusCode"`
	Request    *RequestMeta      `json:"request,omitempty"`
	Message    string            `json:"message,omitempty"`
	Data       T                 `json:"data,omitempty"`
	Error      any               `json:"error,omitempty"` // string | map | struct — only set on failure
}

// RequestMeta captures interesting request context (useful for debugging / audit)
type RequestMeta struct {
	IP           string `json:"ip,omitempty"`
	Method       string `json:"method"`
	Path         string `json:"path"` // cleaner than full URL in most cases
	CorrelationID string `json:"correlationId,omitempty"`
}

// New creates a successful response (most common case)
func NewResponse[T any](status int, message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Success:    true,
		StatusCode: status,
		Message:    message,
		Data:       data,
	}
}

// NewError creates an error-shaped response
func NewError[T any](status int, message string, errDetail any) APIResponse[T] {
	return APIResponse[T]{
		Success:    false,
		StatusCode: status,
		Message:    message,
		Error:      errDetail, // can be string, map[string]any, custom error struct, etc.
	}
}

// WithRequestInfo adds request metadata (call this last, usually in middleware/handler)
func (r *APIResponse[T]) WithRequestInfo(req *http.Request, correlationID string) *APIResponse[T] {
	meta := &RequestMeta{
		Method:       req.Method,
		Path:         req.URL.Path,
		CorrelationID: correlationID,
	}

	// IP logic — you can make this more sophisticated (X-Forwarded-For, etc.)
	if req.RemoteAddr != "" {
		// RemoteAddr is "ip:port" — take only IP
		if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			meta.IP = ip
		} else {
			meta.IP = req.RemoteAddr
		}
	}

	// Hide IP in production (your original logic)
	if strings.ToLower(os.Getenv("GIN_MODE")) == EnvProduction ||
		strings.ToLower(os.Getenv("APP_ENV")) == EnvProduction ||
		os.Getenv("NODE_ENV") == EnvProduction { // ← added for people migrating from Node
		meta.IP = "" // or "[redacted]"
	}

	// Optional: hide correlation ID too (uncomment if desired)
	// if strings.ToLower(os.Getenv("GIN_MODE")) == EnvProduction {
	// 	meta.CorrelationID = ""
	// }

	r.Request = meta
	return r
}
// func CreateUser(c *gin.Context) {
// 	// ... validation failed example
// 	resp := response.NewError[any](
// 		http.StatusBadRequest,
// 		"Validation failed",
// 		map[string]string{"email": "invalid format"},
// 	).WithRequestInfo(c.Request, c.GetString("correlationID"))

// 	c.JSON(resp.StatusCode, resp)
// }
