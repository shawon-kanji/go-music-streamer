package song

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type SongUseCase interface {
	ListSongs(page int, limit int) (*dto.PaginatedSongResponse, error)
	FetchSongByID(id uint) (*dto.SongResponse, error)
}

type songUseCase struct {
	repo repository.SongRepository
}

func NewSongUseCase(repo repository.SongRepository) SongUseCase {
	return &songUseCase{
		repo: repo,
	}
}

func (useCase *songUseCase) ListSongs(page int, limit int) (*dto.PaginatedSongResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	songEntities, totalCount, err := useCase.repo.ListSongs(page, limit)
	if err != nil {
		return nil, err
	}

	var songResponses []*dto.SongResponse
	for _, song := range songEntities {
		songResponses = append(songResponses, &dto.SongResponse{
			ID:        song.ID,
			Title:     song.Title,
			Artist:    song.Artist,
			Album:     song.Album,
			Genre:     song.Genre,
			URL:       song.Url,
			Duration:  song.Duration,
			LikeCount: song.LikeCount,
			Thumbnail: song.Thumbnail,
		})
	}

	// Calculate if there are more
	hasMore := int64(page*limit) < totalCount

	return &dto.PaginatedSongResponse{
		MediaList:  songResponses,
		HasMore:    hasMore,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
	}, nil
}

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
		URL:       song.Url,
		Duration:  song.Duration,
		LikeCount: song.LikeCount,
		Thumbnail: song.Thumbnail,
	}, nil
}
