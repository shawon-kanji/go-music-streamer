package framework

import (
	"errors"
	"fmt"

	"go-music-streamer/internal/domain/apperror"

	"github.com/gin-gonic/gin"
)

// ErrorResponse standardizes error API responses
type ErrorResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"errorCode,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
}

// SendError formats and sends an error JSON response
func SendError(c *gin.Context, statusCode int, message string, errs interface{}) {
	fmt.Println("Error:", message, "Details:", errs) // Log the error details for debugging

	errorCode := ""
	if err, ok := errs.(error); ok {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			errorCode = appErr.Code
			// If we want to hide the full error string and only show exactly the AppError message
			message = appErr.Message
			if len(appErr.Errors) > 0 {
				errs = appErr.Errors // Provide the nested validation errors array
			} else {
				errs = nil // we already send the error code and message clearly
			}
		} else {
			fmt.Println(err, ok)
			errs = err.Error()
		}
	}

	c.JSON(statusCode, ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Errors:    errs,
	})
}

func NotFound(c *gin.Context) {
	SendError(c, 404, "Resource not found", nil)
}

func InternalServerError(c *gin.Context, err error) {
	SendError(c, 500, "Internal server error", err)
}

func BadRequest(c *gin.Context, err error) {
	SendError(c, 400, "Bad request", err)
}

func Unauthorized(c *gin.Context, err error) {
	SendError(c, 401, "Unauthorized", err)
}

func Forbidden(c *gin.Context, err error) {
	SendError(c, 403, "Forbidden", err)
}
