package database

import (
	"go-music-streamer/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DB_HOST + " user=" + cfg.DB_USER + " password=" + cfg.DB_PASSWORD + " dbname=" + cfg.DB_NAME + " port=" + cfg.DB_PORT + " sslmode=" + cfg.DB_SSLMODE
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
