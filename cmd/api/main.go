package main

import (
	"go-music-streamer/internal/api/router"
	"go-music-streamer/internal/config"
	"go-music-streamer/internal/database"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established.", db)

	r := router.New()

	log.Println("Server starting on port " + cfg.APP_PORT + "...")
	if err := r.Run(":" + cfg.APP_PORT); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
