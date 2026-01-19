package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeConflict       ErrorType = "CONFLICT"
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeUnprocessable  ErrorType = "UNPROCESSABLE_ENTITY"
	ErrorTypeRateLimited    ErrorType = "RATE_LIMITED"
)

// AppError represents an application error with context
type AppError struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// ErrorResponse represents the JSON error response
type ErrorResponse struct {
	Type    ErrorType              `json:"type"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ToErrorResponse converts AppError to ErrorResponse
func (e *AppError) ToErrorResponse() ErrorResponse {
	return ErrorResponse{
		Type:    e.Type,
		Message: e.Message,
		Details: e.Details,
	}
}

// New creates a new AppError
func New(errType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		StatusCode: getStatusCode(errType),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		StatusCode: getStatusCode(errType),
		Err:        err,
	}
}

// WithDetails adds details to an error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// getStatusCode returns the HTTP status code for an error type
func getStatusCode(errType ErrorType) int {
	switch errType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeUnprocessable:
		return http.StatusUnprocessableEntity
	case ErrorTypeRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// Common errors
var (
	ErrNotFound         = New(ErrorTypeNotFound, "Resource not found")
	ErrUnauthorized     = New(ErrorTypeUnauthorized, "Unauthorized")
	ErrForbidden        = New(ErrorTypeForbidden, "Forbidden")
	ErrInternal         = New(ErrorTypeInternal, "Internal server error")
	ErrBadRequest       = New(ErrorTypeBadRequest, "Bad request")
	ErrValidation       = New(ErrorTypeValidation, "Validation failed")
	ErrConflict         = New(ErrorTypeConflict, "Resource conflict")
)

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return New(ErrorTypeValidation, message).WithDetails(details)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *AppError {
	return New(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource))
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return New(ErrorTypeConflict, message)
}

// NewInternalError creates a new internal error
func NewInternalError(message string, err error) *AppError {
	return Wrap(err, ErrorTypeInternal, message)
}
