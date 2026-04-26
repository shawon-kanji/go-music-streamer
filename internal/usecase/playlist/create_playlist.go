package playlist

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
)

func (uc *useCase) CreatePlaylist(req *dto.CreatePlaylistRequest) (*dto.PlaylistResponse, error) {
	playlistEntity := &entity.Playlist{
		Name:        req.Name,
		Description: req.Description,
		Visibility:  req.Visibility,
		CreatedBy:   req.CreatedBy,
		Thumbnail:   req.Thumbnail,
	}

	createdPlaylist, err := uc.playlistRepo.CreatePlaylist(playlistEntity)
	if err != nil {
		return nil, err
	}

	return &dto.PlaylistResponse{
		ID:          createdPlaylist.ID,
		Name:        createdPlaylist.Name,
		Description: createdPlaylist.Description,
		Visibility:  createdPlaylist.Visibility,
		CreatedBy:   createdPlaylist.CreatedBy,
		UserDetails: dto.PlaylistUser{
			ID:   createdPlaylist.UserDetails.ID,
			Name: createdPlaylist.UserDetails.Name,
		},
		Thumbnail: createdPlaylist.Thumbnail,
	}, nil
}
