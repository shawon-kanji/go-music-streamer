package repository

import (
	"go-music-streamer/internal/database/postgres"
	"go-music-streamer/internal/domain/entity"

	"gorm.io/gorm"
)

type SongRepository interface {
	ListSongs(page int, limit int) ([]*entity.Song, int64, error)
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) ListSongs(page int, limit int) ([]*entity.Song, int64, error) {
	var songs []*postgres.Song
	var totalCount int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // default limit
	}

	offset := (page - 1) * limit

	// Get total count of songs
	if err := r.db.Model(&postgres.Song{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated songs
	if err := r.db.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		return nil, 0, err
	}

	var songEntities []*entity.Song
	for _, song := range songs {
		songEntities = append(songEntities, &entity.Song{
			ID:        song.ID,
			Title:     song.Title,
			Artist:    song.Artist,
			Album:     song.Album,
			Genre:     song.Genre,
			Duration:  song.Duration,
			Url:       song.Url,
			LikeCount: song.LikeCount,
			Thumbnail: song.Thumbnail,
		})
	}
	return songEntities, totalCount, nil
}
