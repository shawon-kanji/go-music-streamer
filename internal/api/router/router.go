package router

import (
	"go-music-streamer/internal/api/handlers"
	"go-music-streamer/internal/api/middleware"
	"go-music-streamer/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

func New(userHandler *handlers.UserHandler) *gin.Engine {
	r := gin.Default()

	// Fix: [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe.
	// usage: r.SetTrustedProxies([]string{"127.0.0.1"})
	r.SetTrustedProxies(nil)

	r.GET("/health", handlers.Health)

	// User routes with validation middleware
	r.POST("/signup", middleware.ValidateJSON(dto.CreateUserRequest{}, "validatedRequest"), userHandler.Signup)
	r.POST("/login", middleware.ValidateJSON(dto.LoginRequest{}, "validatedRequest"), userHandler.Login)

	return r
}
