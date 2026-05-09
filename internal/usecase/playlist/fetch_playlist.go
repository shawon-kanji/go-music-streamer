package playlist

import (
	"errors"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"

	"gorm.io/gorm"
)

func (uc *useCase) FetchPlaylist(id uint) (*dto.PlaylistResponse, error) {
	playlistEnt, err := uc.playlistRepo.GetPlaylistByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.NotFound, "playlist not found")
		}
		return nil, err
	}

	songsDto := make([]dto.SongResponse, 0, len(playlistEnt.Songs))
	for _, song := range playlistEnt.Songs {
		songsDto = append(songsDto, dto.SongResponse{
			ID:        song.ID,
			Title:     song.Title,
			Artist:    song.Artist,
			Album:     song.Album,
			Duration:  song.Duration,
			URL:       song.URL,
			LikeCount: song.LikeCount,
			Genre:     song.Genre,
			Thumbnail: song.Thumbnail,
		})
	}

	return &dto.PlaylistResponse{
		ID:          playlistEnt.ID,
		Name:        playlistEnt.Name,
		Description: playlistEnt.Description,
		Visibility:  playlistEnt.Visibility,
		CreatedBy:   playlistEnt.CreatedBy,
		UserDetails: dto.PlaylistUser{
			ID:   playlistEnt.UserDetails.ID,
			Name: playlistEnt.UserDetails.Name,
		},
		Thumbnail: playlistEnt.Thumbnail,
		Songs:     songsDto,
	}, nil
}
