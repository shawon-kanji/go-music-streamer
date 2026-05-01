package playlist

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type UseCase interface {
	AddSong(playlistID uint, userID uint, req *dto.AddSongToPlaylistRequest) error
	CreatePlaylist(req *dto.CreatePlaylistRequest) (*dto.PlaylistResponse, error)
	UpdatePlaylist(playlistID uint, userID uint, req *dto.UpdatePlaylistRequest) (*dto.PlaylistResponse, error)
	FetchPlaylist(id uint) (*dto.PlaylistResponse, error)
	FetchPlaylists(page int, limit int, options dto.PlaylistQueryOptions) (*dto.PaginatedPlaylistResponse, error)
}

type useCase struct {
	playlistRepo repository.PlaylistRepository
}

func NewUseCase(playlistRepo repository.PlaylistRepository) UseCase {
	return &useCase{
		playlistRepo: playlistRepo,
	}
}
