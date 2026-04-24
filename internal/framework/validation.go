package framework

import (
	"errors"
	"fmt"

	"go-music-streamer/internal/domain/apperror"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError converts raw gin/validator errors into user-friendly string messages
func FormatValidationError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]string, len(ve))
		for i, fe := range ve {
			out[i] = getErrorMsg(fe)
		}
		appErr := apperror.New("DATA_VALIDATION_ERROR", "Validation failed")
		return appErr.WithErrors(out)
	}
	return apperror.New("DATA_VALIDATION_ERROR", err.Error())
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email format", fe.Field())
	case "min":
		return fmt.Sprintf("%s length must be at least %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s length must be at most %s", fe.Field(), fe.Param())
	case "alpha":
		return fmt.Sprintf("%s can only contain alphabetic characters", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s can only contain alphanumeric characters", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s can only contain numeric characters", fe.Field())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	}
	// Fallback error message
	return fmt.Sprintf("%s is invalid", fe.Field())
}
