package repository

import (
	"go-music-streamer/internal/database/postgres"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"

	"gorm.io/gorm"
)

type PlaylistRepository interface {
	// Other playlist-related methods...
	CreatePlaylist(playlist *entity.Playlist) (entity.Playlist, error)
	GetPlaylistByID(id uint) (*entity.Playlist, error)
	UpdatePlaylist(playlist *entity.Playlist) (entity.Playlist, error)
	DeletePlaylist(id uint) error
	ListPlaylists(page int, limit int, options dto.PlaylistQueryOptions) ([]*entity.Playlist, int64, error)
	AddSongToPlaylist(playlistID uint, songID uint) error
}

type playlistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) CreatePlaylist(playlist *entity.Playlist) (entity.Playlist, error) {
	playlistModel := &postgres.Playlist{
		Name:        playlist.Name,
		Description: playlist.Description,
		Visibility:  playlist.Visibility,
		CreatedBy:   playlist.CreatedBy,
		Thumbnail:   playlist.Thumbnail,
	}

	if err := r.db.Create(playlistModel).Error; err != nil {
		return entity.Playlist{}, err
	}
	playlist.ID = playlistModel.ID
	return *playlist, nil
}

func (r *playlistRepository) GetPlaylistByID(id uint) (*entity.Playlist, error) {
	var playlistModel postgres.Playlist

	if err := r.db.Preload("UserDetails").Preload("Songs").First(&playlistModel, id).Error; err != nil {
		return nil, err
	}

	playlist := &entity.Playlist{
		ID:          playlistModel.ID,
		Name:        playlistModel.Name,
		Description: playlistModel.Description,
		Visibility:  playlistModel.Visibility,
		CreatedBy:   playlistModel.CreatedBy,
		UserDetails: entity.User{
			ID:    playlistModel.UserDetails.ID,
			Name:  playlistModel.UserDetails.Username,
			Email: playlistModel.UserDetails.Email,
		},
		Thumbnail: playlistModel.Thumbnail,
	}

	for _, songModel := range playlistModel.Songs {
		playlist.Songs = append(playlist.Songs, entity.Song{
			ID:        songModel.ID,
			Title:     songModel.Title,
			Artist:    songModel.Artist,
			Album:     songModel.Album,
			Genre:     songModel.Genre,
			Duration:  songModel.Duration,
			Url:       songModel.Url,
			LikeCount: songModel.LikeCount,
			Thumbnail: songModel.Thumbnail,
		})
	}

	return playlist, nil
}

func (r *playlistRepository) UpdatePlaylist(playlist *entity.Playlist) (entity.Playlist, error) {
	var playlistModel postgres.Playlist

	if err := r.db.First(&playlistModel, playlist.ID).Error; err != nil {
		return entity.Playlist{}, err
	}

	playlistModel.Name = playlist.Name
	playlistModel.Description = playlist.Description
	playlistModel.Visibility = playlist.Visibility
	playlistModel.Thumbnail = playlist.Thumbnail

	if err := r.db.Save(&playlistModel).Error; err != nil {
		return entity.Playlist{}, err
	}

	return *playlist, nil
}

func (r *playlistRepository) DeletePlaylist(id uint) error {
	return r.db.Delete(&postgres.Playlist{}, id).Error
}

func (r *playlistRepository) ListPlaylists(page int, limit int, options dto.PlaylistQueryOptions) ([]*entity.Playlist, int64, error) {
	var playlists []*postgres.Playlist
	var totalCount int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // default limit
	}

	offset := (page - 1) * limit

	// Grouping the OR condition ensures that if future AND conditions
	// (or automatic soft-delete scopes) are appended, the precedence remains correct.
	baseQuery := r.db.Model(&postgres.Playlist{}).
		Where(
			r.db.Where("visibility = ?", "public").
				Or("created_by = ?", options.UserID),
		)

	// Get total count of playlists
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated playlists
	if err := baseQuery.Preload("UserDetails").Limit(limit).Offset(offset).Find(&playlists).Error; err != nil {
		return nil, 0, err
	}
	var playlistEntities []*entity.Playlist
	// no need to preload songs here, as they will be fetched separately when fetching playlist details
	for _, playlist := range playlists {
		playlistEntities = append(playlistEntities, &entity.Playlist{
			ID:          playlist.ID,
			Name:        playlist.Name,
			Description: playlist.Description,
			Visibility:  playlist.Visibility,
			CreatedBy:   playlist.CreatedBy,
			UserDetails: entity.User{
				ID:    playlist.UserDetails.ID,
				Name:  playlist.UserDetails.Username,
				Email: playlist.UserDetails.Email,
			},
			Thumbnail: playlist.Thumbnail,
			Songs:     make([]entity.Song, 0), // Songs will be fetched separately in GetPlaylistByID
		})
	}
	return playlistEntities, totalCount, nil
}

func (r *playlistRepository) AddSongToPlaylist(playlistID uint, songID uint) error {
	playlist := &postgres.Playlist{}
	if err := r.db.First(playlist, playlistID).Error; err != nil {
		return err
	}
	song := &postgres.Song{}
	if err := r.db.First(song, songID).Error; err != nil {
		return err
	}
	return r.db.Model(playlist).Association("Songs").Append(song)
}
