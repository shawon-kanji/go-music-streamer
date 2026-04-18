package postgres

import (
	"go-music-streamer/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.Postgres.Host + " user=" + cfg.Postgres.User + " password=" + cfg.Postgres.Password + " dbname=" + cfg.Postgres.Name + " port=" + cfg.Postgres.Port + " sslmode=" + cfg.Postgres.SSLMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
