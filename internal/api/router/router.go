package router

import (
	"go-music-streamer/internal/api/handlers"
	"go-music-streamer/internal/api/middleware"
	"go-music-streamer/internal/bootstrap"
	"go-music-streamer/internal/domain/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// New initializes the Gin application engine and configures all routing groups and middleware.
func New(db *gorm.DB, appHandlers *bootstrap.AppHandlers) *gin.Engine {
	r := gin.Default()

	r.MaxMultipartMemory = 10 << 20 // 10 MiB

	// Fix: [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe.
	r.SetTrustedProxies(nil)

	// Public routes
	r.GET("/health", handlers.Health)

	publicUserRoutes := r.Group("/auth")
	{
		publicUserRoutes.POST("/signup", middleware.ValidateJSON(dto.CreateUserRequest{}, "validatedRequest"), appHandlers.UserSignupHandler.Handle)
		publicUserRoutes.POST("/login", middleware.ValidateJSON(dto.LoginRequest{}, "validatedRequest"), appHandlers.UserLoginHandler.Handle)
	}

	// Secured routes requiring JWT authentication
	secureRoutes := r.Group("/api/v1")
	secureRoutes.Use(middleware.Authenticate())
	{
		// Any authenticated user can check their profile
		secureRoutes.GET("/me", appHandlers.UserProfileHandler.Handle)
		secureRoutes.GET("/songs", appHandlers.ListSongsHandler.Handle)
		secureRoutes.GET("/songs/:id", appHandlers.FetchSongHandler.Handle)
		secureRoutes.PATCH("/songs/:id", middleware.ValidateJSON(dto.UpdateSongRequest{}, "validatedRequest"), appHandlers.UpdateSongHandler.Handle)
		secureRoutes.POST("/songs/upload", appHandlers.SongHandler.UploadSong)

		// Playlist routes
		secureRoutes.POST("/playlists", middleware.ValidateJSON(dto.CreatePlaylistRequest{}, "validatedRequest"), appHandlers.CreatePlaylistHandler.Handle)
		secureRoutes.PATCH("/playlists/:id", middleware.ValidateJSON(dto.UpdatePlaylistRequest{}, "validatedRequest"), appHandlers.UpdatePlaylistHandler.Handle)
		secureRoutes.GET("/playlists", appHandlers.FetchPlaylistsHandler.Handle)
		secureRoutes.GET("/playlists/:id", appHandlers.FetchPlaylistHandler.Handle)
		secureRoutes.POST("/playlists/:id/songs", middleware.ValidateJSON(dto.AddSongToPlaylistRequest{}, "validatedRequest"), appHandlers.AddSongToPlaylistHandler.Handle)

		// Example: Only users with the explicit "READ" permission on "USER" resource can access this specific group.
		adminRoutes := secureRoutes.Group("/admin")
		adminRoutes.Use(middleware.Authorize(db, "USER", "READ"))
		{
			// adminRoutes.GET("/users", adminHandler.ListUsers)
		}
	}

	return r
}
