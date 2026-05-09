package repository

import (
	"errors"

	"go-music-streamer/internal/database/postgres"
	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/entity"

	"gorm.io/gorm"
)

type SongRepository interface {
	CreateSong(song *entity.Song) (*entity.Song, error)
	ListSongs(page int, limit int) ([]*entity.Song, int64, error)
	GetSongByID(id uint) (*entity.Song, error)
	UpdateSong(song *entity.Song) (*entity.Song, error)
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) CreateSong(song *entity.Song) (*entity.Song, error) {
	songModel := &postgres.Song{
		Title:     song.Title,
		Artist:    song.Artist,
		Album:     song.Album,
		Genre:     song.Genre,
		Duration:  song.Duration,
		URL:       song.URL,
		LikeCount: song.LikeCount,
		Thumbnail: song.Thumbnail,
	}

	if err := r.db.Create(songModel).Error; err != nil {
		return nil, err
	}

	return &entity.Song{
		ID:        songModel.ID,
		Title:     songModel.Title,
		Artist:    songModel.Artist,
		Album:     songModel.Album,
		Genre:     songModel.Genre,
		Duration:  songModel.Duration,
		URL:       songModel.URL,
		LikeCount: songModel.LikeCount,
		Thumbnail: songModel.Thumbnail,
	}, nil
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
			URL:       song.URL,
			LikeCount: song.LikeCount,
			Thumbnail: song.Thumbnail,
		})
	}
	return songEntities, totalCount, nil
}

func (r *songRepository) GetSongByID(id uint) (*entity.Song, error) {
	var song postgres.Song
	err := r.db.First(&song, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.NotFound, "song not found")
		}
		return nil, err
	}

	return &entity.Song{
		ID:        song.ID,
		Title:     song.Title,
		Artist:    song.Artist,
		Album:     song.Album,
		Genre:     song.Genre,
		Duration:  song.Duration,
		URL:       song.URL,
		LikeCount: song.LikeCount,
		Thumbnail: song.Thumbnail,
	}, nil
}

func (r *songRepository) UpdateSong(song *entity.Song) (*entity.Song, error) {
	var songModel postgres.Song
	if err := r.db.First(&songModel, song.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.NotFound, "song not found")
		}
		return nil, err
	}

	songModel.Title = song.Title
	songModel.Artist = song.Artist
	songModel.Album = song.Album
	songModel.Genre = song.Genre
	songModel.Duration = song.Duration
	songModel.URL = song.URL
	songModel.Thumbnail = song.Thumbnail

	if err := r.db.Save(&songModel).Error; err != nil {
		return nil, err
	}

	return &entity.Song{
		ID:        songModel.ID,
		Title:     songModel.Title,
		Artist:    songModel.Artist,
		Album:     songModel.Album,
		Genre:     songModel.Genre,
		Duration:  songModel.Duration,
		URL:       songModel.URL,
		LikeCount: songModel.LikeCount,
		Thumbnail: songModel.Thumbnail,
	}, nil
}
