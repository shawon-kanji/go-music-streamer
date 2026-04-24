package apperror

import "fmt"

// AppError is a custom error type that includes a unique error code
// and a human-readable message.
type AppError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"` // Optional field for additional error details
}

// Error implements the standard Go error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New creates a new AppError with the given code and message.
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new AppError with formatting for the message.
func Newf(code, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// WithErrors adds detailed sub-errors to the AppError
func (e *AppError) WithErrors(errs []string) *AppError {
	e.Errors = errs
	return e
}
