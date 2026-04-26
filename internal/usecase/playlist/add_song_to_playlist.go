package playlist

import (
	"go-music-streamer/internal/domain/dto"
)

func (uc *useCase) AddSong(playlistID uint, req *dto.AddSongToPlaylistRequest) error {
	return uc.playlistRepo.AddSongToPlaylist(playlistID, req.SongID)
}
