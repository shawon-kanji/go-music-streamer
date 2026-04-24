package framework

import (
	"github.com/gin-gonic/gin"
)

// Response standardizes successful API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendSuccess formats and sends a successful JSON response
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}
