package song

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
)

func (useCase *songUseCase) CreateSong(req *dto.CreateSongRequest) (*dto.SongResponse, error) {
	createdSong, err := useCase.repo.CreateSong(&entity.Song{
		Title:     req.Title,
		Artist:    req.Artist,
		Album:     req.Album,
		Genre:     req.Genre,
		Duration:  req.Duration,
		URL:       req.URL,
		Thumbnail: req.Thumbnail,
	})
	if err != nil {
		return nil, err
	}

	return &dto.SongResponse{
		ID:        createdSong.ID,
		Title:     createdSong.Title,
		Artist:    createdSong.Artist,
		Album:     createdSong.Album,
		Genre:     createdSong.Genre,
		Duration:  createdSong.Duration,
		URL:       createdSong.URL,
		LikeCount: createdSong.LikeCount,
		Thumbnail: createdSong.Thumbnail,
	}, nil
}