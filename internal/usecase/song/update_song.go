package song

import (
	"reflect"

	"go-music-streamer/internal/domain/dto"
)

func (useCase *songUseCase) UpdateSong(id uint, req *dto.UpdateSongRequest) (*dto.SongResponse, error) {
	song, err := useCase.repo.GetSongByID(id)
	if err != nil {
		return nil, err
	}

	applySongUpdates(song, req)

	updatedSong, err := useCase.repo.UpdateSong(song)
	if err != nil {
		return nil, err
	}

	return &dto.SongResponse{
		ID:        updatedSong.ID,
		Title:     updatedSong.Title,
		Artist:    updatedSong.Artist,
		Album:     updatedSong.Album,
		Genre:     updatedSong.Genre,
		URL:       updatedSong.Url,
		Duration:  updatedSong.Duration,
		LikeCount: updatedSong.LikeCount,
		Thumbnail: updatedSong.Thumbnail,
	}, nil
}

// applySongUpdates copies non-nil pointer fields from request DTO into matching entity fields.
func applySongUpdates(target interface{}, req *dto.UpdateSongRequest) {
	fieldMapping := map[string]string{
		"URL": "Url",
	}

	reqVal := reflect.ValueOf(req).Elem()
	reqType := reqVal.Type()
	targetVal := reflect.ValueOf(target).Elem()

	for i := 0; i < reqVal.NumField(); i++ {
		fieldVal := reqVal.Field(i)
		if fieldVal.IsNil() {
			continue
		}

		reqFieldName := reqType.Field(i).Name
		targetFieldName, ok := fieldMapping[reqFieldName]
		if !ok {
			targetFieldName = reqFieldName
		}

		targetField := targetVal.FieldByName(targetFieldName)
		if !targetField.IsValid() || !targetField.CanSet() || targetField.Type() != fieldVal.Elem().Type() {
			continue
		}

		targetField.Set(fieldVal.Elem())
	}
}
