package middleware

import (
	"reflect"

	"go-music-streamer/internal/framework"

	"github.com/gin-gonic/gin"
)

// ValidateJSON creates a middleware that parses and validates incoming JSON against the given struct type.
// Upon success, it stores the parsed DTO pointer in the Gin context under contextKey.
func ValidateJSON(dtoType interface{}, contextKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if contextKey == "" {
			contextKey = "body"
		}
		// reflect.TypeOf gets the struct type
		t := reflect.TypeOf(dtoType)
		if t.Kind() == reflect.Ptr {
			t = t.Elem() // Dereference if a pointer was passed
		}

		// Create a new empty instance of the struct pointer
		v := reflect.New(t).Interface()

		if err := c.ShouldBindJSON(v); err != nil {
			errs := framework.FormatValidationError(err)
			framework.BadRequest(c, errs)
			c.Abort() // Stop the handler chain
			return
		}

		// Store the populated struct pointer in the context for down-stream handlers
		c.Set(contextKey, v)
		c.Next()
	}
}
