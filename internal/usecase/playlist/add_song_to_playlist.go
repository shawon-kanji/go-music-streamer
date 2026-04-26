package playlist

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type AddSongToPlaylist interface {
	AddSong(playlistID uint, req *dto.AddSongToPlaylistRequest) error
}

type addSongToPlaylist struct {
	playlistRepo repository.PlaylistRepository
}

func NewAddSongToPlaylist(playlistRepo repository.PlaylistRepository) AddSongToPlaylist {
	return &addSongToPlaylist{
		playlistRepo: playlistRepo,
	}
}

func (uc *addSongToPlaylist) AddSong(playlistID uint, req *dto.AddSongToPlaylistRequest) error {
	return uc.playlistRepo.AddSongToPlaylist(playlistID, req.SongID)
}
