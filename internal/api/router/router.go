package router

import (
	"go-music-streamer/internal/api/handlers"
	"go-music-streamer/internal/api/middleware"
	"go-music-streamer/internal/domain/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// New initializes the Gin application engine and configures all routing groups and middleware.
func New(db *gorm.DB, userHandler *handlers.UserHandler) *gin.Engine {
	r := gin.Default()

	// Fix: [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe.
	r.SetTrustedProxies(nil)

	// Public routes
	r.GET("/health", handlers.Health)

	publicUserRoutes := r.Group("/auth")
	{
		publicUserRoutes.POST("/signup", middleware.ValidateJSON(dto.CreateUserRequest{}, "validatedRequest"), userHandler.Signup)
		publicUserRoutes.POST("/login", middleware.ValidateJSON(dto.LoginRequest{}, "validatedRequest"), userHandler.Login)
	}

	// Secured routes requiring JWT authentication
	secureRoutes := r.Group("/api/v1")
	secureRoutes.Use(middleware.Authenticate())
	{
		// Any authenticated user can check their profile
		secureRoutes.GET("/me", userHandler.Profile)
		secureRoutes.GET("/songs", handlers.NewSongHandler(nil).ListSongs) // Placeholder for SongHandler with actual use case

		// Example: Only users with the explicit "READ" permission on "USER" resource can access this specific group.
		adminRoutes := secureRoutes.Group("/admin")
		adminRoutes.Use(middleware.Authorize(db, "USER", "READ"))
		{
			// adminRoutes.GET("/users", adminHandler.ListUsers)
		}
	}

	return r
}
