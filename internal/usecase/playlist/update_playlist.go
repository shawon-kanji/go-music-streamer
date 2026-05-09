package playlist

import (
	"errors"
	"reflect"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"

	"gorm.io/gorm"
)

func (uc *useCase) UpdatePlaylist(playlistID uint, userID uint, req *dto.UpdatePlaylistRequest) (*dto.PlaylistResponse, error) {
	playlistEnt, err := uc.playlistRepo.GetPlaylistByID(playlistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.NotFound, "playlist not found")
		}
		return nil, err
	}

	if playlistEnt.CreatedBy != userID {
		return nil, apperror.New(apperror.Unauthorized, "only the playlist creator can edit this playlist")
	}

	applyPlaylistUpdates(playlistEnt, req)

	updatedPlaylist, err := uc.playlistRepo.UpdatePlaylist(playlistEnt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.NotFound, "playlist not found")
		}
		return nil, err
	}

	songs := make([]dto.SongResponse, 0, len(updatedPlaylist.Songs))
	for _, song := range updatedPlaylist.Songs {
		songs = append(songs, dto.SongResponse{
			ID:        song.ID,
			Title:     song.Title,
			Artist:    song.Artist,
			Album:     song.Album,
			Genre:     song.Genre,
			URL:       song.URL,
			Duration:  song.Duration,
			LikeCount: song.LikeCount,
			Thumbnail: song.Thumbnail,
		})
	}

	return &dto.PlaylistResponse{
		ID:          updatedPlaylist.ID,
		Name:        updatedPlaylist.Name,
		Description: updatedPlaylist.Description,
		Visibility:  updatedPlaylist.Visibility,
		CreatedBy:   updatedPlaylist.CreatedBy,
		UserDetails: dto.PlaylistUser{
			ID:   updatedPlaylist.UserDetails.ID,
			Name: updatedPlaylist.UserDetails.Name,
		},
		Thumbnail: updatedPlaylist.Thumbnail,
		Songs:     songs,
	}, nil
}

// applyPlaylistUpdates copies non-empty fields from request DTO into entity fields with matching names.
func applyPlaylistUpdates(target interface{}, req *dto.UpdatePlaylistRequest) {
	reqVal := reflect.ValueOf(req).Elem()
	reqType := reqVal.Type()
	targetVal := reflect.ValueOf(target).Elem()

	for i := 0; i < reqVal.NumField(); i++ {
		fieldVal := reqVal.Field(i)

		if fieldVal.Kind() == reflect.String && fieldVal.String() == "" {
			continue
		}

		targetField := targetVal.FieldByName(reqType.Field(i).Name)
		if !targetField.IsValid() || !targetField.CanSet() || targetField.Type() != fieldVal.Type() {
			continue
		}

		targetField.Set(fieldVal)
	}
}
