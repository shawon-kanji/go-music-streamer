package playlist

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type FetchPlaylists interface {
	FetchPlaylists(page int, limit int) (*dto.PaginatedPlaylistResponse, error)
}

type fetchPlaylists struct {
	playlistRepo repository.PlaylistRepository
}

func NewFetchPlaylists(playlistRepo repository.PlaylistRepository) FetchPlaylists {
	return &fetchPlaylists{
		playlistRepo: playlistRepo,
	}
}

func (uc *fetchPlaylists) FetchPlaylists(page int, limit int) (*dto.PaginatedPlaylistResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	playlists, totalCount, err := uc.playlistRepo.ListPlaylists(page, limit)
	if err != nil {
		return nil, err
	}

	var playlistResponses []*dto.PlaylistResponse
	for _, p := range playlists {
		playlistResponses = append(playlistResponses, &dto.PlaylistResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Visibility:  p.Visibility,
			CreatedBy:   p.CreatedBy,
			Thumbnail:   p.Thumbnail,
			UserDetails: dto.PlaylistUser{
				ID:   p.UserDetails.ID,
				Name: p.UserDetails.Name,
			},
		})
	}

	hasMore := int64(page*limit) < totalCount

	return &dto.PaginatedPlaylistResponse{
		PlaylistList: playlistResponses,
		HasMore:      hasMore,
		Page:         page,
		Limit:        limit,
		TotalCount:   totalCount,
	}, nil
}
