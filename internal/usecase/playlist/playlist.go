package playlist

import (
	"errors"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
	"go-music-streamer/internal/repository"
	"gorm.io/gorm"
)

type UseCase interface {
	AddSong(playlistID uint, req *dto.AddSongToPlaylistRequest) error
	CreatePlaylist(req *dto.CreatePlaylistRequest) (*dto.PlaylistResponse, error)
	FetchPlaylist(id uint) (*dto.PlaylistResponse, error)
	FetchPlaylists(page int, limit int) (*dto.PaginatedPlaylistResponse, error)
}

type useCase struct {
	playlistRepo repository.PlaylistRepository
}

func NewUseCase(playlistRepo repository.PlaylistRepository) UseCase {
	return &useCase{
		playlistRepo: playlistRepo,
	}
}

func (uc *useCase) AddSong(playlistID uint, req *dto.AddSongToPlaylistRequest) error {
	return uc.playlistRepo.AddSongToPlaylist(playlistID, req.SongID)
}

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
			URL:       song.Url,
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
		Thumbnail:   playlistEnt.Thumbnail,
		Songs:       songsDto,
	}, nil
}

func (uc *useCase) FetchPlaylists(page int, limit int) (*dto.PaginatedPlaylistResponse, error) {
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
