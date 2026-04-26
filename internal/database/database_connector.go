package database

import (
	"go-music-streamer/internal/config"
	"go-music-streamer/internal/database/postgres"
	"log"
)

type DbClient struct {
	clients map[string]interface{}
}

func ConnectDBSources(cfg config.Config) (*DbClient, error) {
	availableSources := cfg.DBSources
	dbClient := &DbClient{
		clients: make(map[string]interface{}),
	}

	for _, source := range availableSources {
		switch source {
		case "postgres":
			db, err := postgres.ConnectDB(cfg)
			if err != nil {
				return nil, err
			}
			dbClient.clients["postgres"] = db
			err = db.AutoMigrate(
				&postgres.User{},
				&postgres.UserRole{},
				&postgres.Role{},
				&postgres.RolePermission{},
				&postgres.Resource{},
				&postgres.Action{},
				&postgres.Admin{},
				&postgres.Playlist{},
				&postgres.Song{},
				&postgres.UserLikedSong{},
				&postgres.AdminRole{},
			)

			if err != nil {
				log.Fatalf("Failed to apply database migrations: %v", err)
			}

			log.Println("Database schemas migrated successfully.")
			// Add more sources here as needed
		default:
			// ignore unknown sources for now
		}
	}
	return dbClient, nil

}

func (c *DbClient) GetPostgres() interface{} {
	return c.clients["postgres"]
}
