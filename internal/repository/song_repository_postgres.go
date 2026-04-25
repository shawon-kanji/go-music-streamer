package repository

import (
	"go-music-streamer/internal/domain/entity"

	"gorm.io/gorm"
)

type SongRepository interface {
	ListSongs() ([]*entity.Song, error)
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) ListSongs() ([]*entity.Song, error) {
	var songs []*entity.Song
	if err := r.db.Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}
