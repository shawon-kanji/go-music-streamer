package song

import (
	"go-music-streamer/internal/domain/dto"
)

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
			URL:       song.URL,
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
