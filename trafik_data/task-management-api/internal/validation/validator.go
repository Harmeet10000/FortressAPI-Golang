package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourusername/task-management-api/internal/errs"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// ValidationErrors holds multiple validation errors
type ValidationErrors map[string][]string

// Add adds a validation error for a field
func (v ValidationErrors) Add(field, message string) {
	if v[field] == nil {
		v[field] = []string{}
	}
	v[field] = append(v[field], message)
}

// HasErrors returns true if there are validation errors
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// ToAppError converts validation errors to AppError
func (v ValidationErrors) ToAppError() *errs.AppError {
	if !v.HasErrors() {
		return nil
	}

	details := make(map[string]interface{})
	for field, errors := range v {
		details[field] = errors
	}

	return errs.NewValidationError("Validation failed", details)
}

// Validator provides validation utilities
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors),
	}
}

// Required validates that a field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors.Add(field, fmt.Sprintf("%s is required", field))
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors.Add(field, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors.Add(field, fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	if value != "" && !emailRegex.MatchString(value) {
		v.errors.Add(field, fmt.Sprintf("%s must be a valid email address", field))
	}
	return v
}

// UUID validates UUID format
func (v *Validator) UUID(field, value string) *Validator {
	if value != "" && !uuidRegex.MatchString(value) {
		v.errors.Add(field, fmt.Sprintf("%s must be a valid UUID", field))
	}
	return v
}

// Min validates minimum numeric value
func (v *Validator) Min(field string, value, min int) *Validator {
	if value < min {
		v.errors.Add(field, fmt.Sprintf("%s must be at least %d", field, min))
	}
	return v
}

// Max validates maximum numeric value
func (v *Validator) Max(field string, value, max int) *Validator {
	if value > max {
		v.errors.Add(field, fmt.Sprintf("%s must be at most %d", field, max))
	}
	return v
}

// In validates that value is in allowed list
func (v *Validator) In(field, value string, allowed []string) *Validator {
	if value == "" {
		return v
	}

	for _, a := range allowed {
		if value == a {
			return v
		}
	}

	v.errors.Add(field, fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowed, ", ")))
	return v
}

// Custom allows custom validation
func (v *Validator) Custom(field string, isValid bool, message string) *Validator {
	if !isValid {
		v.errors.Add(field, message)
	}
	return v
}

// Validate returns validation error if any
func (v *Validator) Validate() error {
	return v.errors.ToAppError()
}

// IsValidEmail checks if a string is a valid email
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}
