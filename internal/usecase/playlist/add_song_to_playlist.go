package playlist

import (
	"errors"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"

	"gorm.io/gorm"
)

func (uc *useCase) AddSong(playlistID uint, userID uint, req *dto.AddSongToPlaylistRequest) error {
	playlistEnt, err := uc.playlistRepo.GetPlaylistByID(playlistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.NotFound, "playlist not found")
		}
		return err
	}

	if playlistEnt.CreatedBy != userID {
		return apperror.New(apperror.Unauthorized, "only the playlist creator can edit this playlist")
	}

	if err := uc.playlistRepo.AddSongToPlaylist(playlistID, req.SongID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.New(apperror.NotFound, "playlist or song not found")
		}
		return err
	}

	return nil
}
