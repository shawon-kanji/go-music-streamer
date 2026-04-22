package main

import (
	"log"

	"go-music-streamer/internal/api/handlers"
	"go-music-streamer/internal/api/router"
	"go-music-streamer/internal/config"
	"go-music-streamer/internal/database"
	"go-music-streamer/internal/repository"
	"go-music-streamer/internal/usecase/user"

	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbClient, err := database.ConnectDBSources(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to databases: %v", err)
	}

	// Get Postgres DB
	pgDB, ok := dbClient.GetPostgres().(*gorm.DB)
	if !ok || pgDB == nil {
		log.Fatalf("Failed to retrieve Postgres connection")
	}
	log.Println("Database connection established.")

	// Wire dependencies: Repository -> UseCase -> Handler
	userRepo := repository.NewUserRepository(pgDB)
	userUseCase := user.NewUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(userUseCase)

	r := router.New(userHandler)

	log.Println("Server starting on port " + cfg.AppPort + "...")
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
