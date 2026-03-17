package router

import (
	"github.com/gin-gonic/gin"
	"go-music-streamer/internal/api/handlers"
)

func New() *gin.Engine {
	r := gin.Default()

	// Fix: [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe.
	// usage: r.SetTrustedProxies([]string{"127.0.0.1"})
	r.SetTrustedProxies(nil)

	r.GET("/health", handlers.Health)

	return r
}
