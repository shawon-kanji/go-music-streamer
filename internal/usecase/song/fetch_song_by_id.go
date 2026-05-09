package song

import (
	"go-music-streamer/internal/domain/dto"
)

func (useCase *songUseCase) FetchSongByID(id uint) (*dto.SongResponse, error) {
	song, err := useCase.repo.GetSongByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.SongResponse{
		ID:        song.ID,
		Title:     song.Title,
		Artist:    song.Artist,
		Album:     song.Album,
		Genre:     song.Genre,
		URL:       song.URL,
		Duration:  song.Duration,
		LikeCount: song.LikeCount,
		Thumbnail: song.Thumbnail,
	}, nil
}
