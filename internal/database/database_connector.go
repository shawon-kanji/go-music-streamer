package database

import (
	"go-music-streamer/internal/config"
	"go-music-streamer/internal/database/postgres"
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
			// Add more sources here as needed
		default:
			return nil, nil
		}
	}
	return nil, nil

}
