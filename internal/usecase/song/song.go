package song

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type SongUseCase interface {
	ListSongs(page int, limit int) (*dto.PaginatedSongResponse, error)
	FetchSongByID(id uint) (*dto.SongResponse, error)
	UpdateSong(id uint, req *dto.UpdateSongRequest) (*dto.SongResponse, error)
}

type songUseCase struct {
	repo repository.SongRepository
}

func NewSongUseCase(repo repository.SongRepository) SongUseCase {
	return &songUseCase{
		repo: repo,
	}
}
