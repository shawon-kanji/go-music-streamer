package song

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
	"go-music-streamer/internal/repository"

	"github.com/gin-gonic/gin"
)

type SongUseCase interface {
	CreateSong(song *dto.CreateSongRequest) (*dto.SongResponse, error)
	ListSongs(page int, limit int) (*dto.PaginatedSongResponse, error)
	FetchSongByID(id uint) (*dto.SongResponse, error)
	UpdateSong(id uint, req *dto.UpdateSongRequest) (*dto.SongResponse, error)
	UploadSong(req *dto.UploadSongRequest, c *gin.Context) (*entity.Song, error)
	SearchSongs(req *dto.SearchSongsRequest) (*dto.SemanticSearchResponse, error)
}

type songUseCase struct {
	repo repository.SongRepository
}

func NewSongUseCase(repo repository.SongRepository) SongUseCase {
	return &songUseCase{
		repo: repo,
	}
}
