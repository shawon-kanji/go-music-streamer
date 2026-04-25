package song

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type SongUseCase interface {
	ListSongs() ([]*dto.SongResponse, error)
}

type songUseCase struct {
	repo repository.SongRepository
}

func NewSongUseCase(repo repository.SongRepository) SongUseCase {
	return &songUseCase{
		repo: repo,
	}
}

func (useCase *songUseCase) ListSongs() ([]*dto.SongResponse, error) {
	songEntities, err := useCase.repo.ListSongs()
	if err != nil {
		return nil, err
	}

	var songResponses []*dto.SongResponse
	for _, song := range songEntities {
		songResponses = append(songResponses, &dto.SongResponse{
			ID:     song.ID,
			Title:  song.Title,
			Artist: song.Artist,
			Album:  song.Album,
			Genre:  song.Genre,
			URL:    song.Url,
		})
	}
	return songResponses, nil
}
