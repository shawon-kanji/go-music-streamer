package middleware

import (
	"fmt"
	"strings"

	"go-music-streamer/internal/config"
	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/framework"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Authenticate verifies the JWT token present in the Authorization header.
// If valid, it extracts the user ID and sets it in the gin Context.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "invalid authorization header format. Expected 'Bearer <token>'"))
			c.Abort()
			return
		}

		tokenStr := parts[1]
		appConfig := config.GetConfig()

		// Parse and validate the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(appConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "invalid token claims"))
			c.Abort()
			return
		}

		userID, ok := claims["sub"]
		if !ok {
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "user identity not found in token"))
			c.Abort()
			return
		}

		var uid uint
		switch v := userID.(type) {
		case float64:
			uid = uint(v)
		case int:
			uid = uint(v)
		case uint:
			uid = v
		default:
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "invalid user identity format"))
			c.Abort()
			return
		}

		// Store the userID directly for the Authorize middleware and downstream handlers
		c.Set("userID", uid)
		c.Next()
	}
}

// Authorize checks if the authenticated user has permission to perform an action on a resource.
// It assumes that an authentication middleware has already set "userID" in the *gin.Context.
func Authorize(db *gorm.DB, requiredResource string, requiredAction string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the authenticated user ID from context.
		// Adjust the key ("userID") and type cast based on your Auth middleware implementation.
		userIDVal, exists := c.Get("userID")
		if !exists {
			framework.Unauthorized(c, apperror.New(apperror.Unauthorized, "user not authenticated"))
			c.Abort()
			return
		}

		var userID uint
		switch v := userIDVal.(type) {
		case uint:
			userID = v
		case float64: // JSON numbers are usually parsed as float64 (JWT payloads)
			userID = uint(v)
		case int:
			userID = uint(v)
		default:
			framework.InternalServerError(c, apperror.New(apperror.InternalError, "invalid user identity inside context"))
			c.Abort()
			return
		}

		// 2. Query the database to check if a valid RolePermission exists.
		// We join user_roles -> roles -> role_permissions -> resources and actions.
		// We handle exact matches as well as the wildcard "*".
		var count int64
		err := db.Table("user_roles").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
			Joins("JOIN resources ON resources.id = role_permissions.resource_id").
			Joins("JOIN actions ON actions.id = role_permissions.action_id").
			Where("user_roles.user_id = ?", userID).
			Where("(resources.name = ? OR resources.name = ?)", requiredResource, "*").
			Where("(actions.name = ? OR actions.name = ?)", requiredAction, "*").
			Where("user_roles.deleted_at IS NULL").
			Where("roles.deleted_at IS NULL").
			Where("role_permissions.deleted_at IS NULL").
			Where("resources.deleted_at IS NULL").
			Where("actions.deleted_at IS NULL").
			Count(&count).Error

		if err != nil {
			framework.InternalServerError(c, apperror.New(apperror.InternalError, "failed to verify permissions"))
			c.Abort()
			return
		}

		// 3. Deny if no valid role permission was found
		if count == 0 {
			framework.Forbidden(c, apperror.New(apperror.Unauthorized, "you do not have permission to perform this action"))
			c.Abort()
			return
		}

		// Proceed to handler
		c.Next()
	}
}
