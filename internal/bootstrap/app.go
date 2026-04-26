package bootstrap

import (
	"go-music-streamer/internal/api/handlers"
	"go-music-streamer/internal/repository"
	"go-music-streamer/internal/usecase/playlist"
	"go-music-streamer/internal/usecase/song"
	"go-music-streamer/internal/usecase/user"

	"gorm.io/gorm"
)

// AppHandlers struct acts as a container for all HTTP handlers
type AppHandlers struct {
	UserHandler     *handlers.UserHandler
	SongHandler     *handlers.SongHandler
	PlaylistHandler *handlers.PlaylistHandler
}

// InitHandlers initializes and injects all repositories, use cases, and handlers.
func InitHandlers(db *gorm.DB) *AppHandlers {
	// 1. Repositories
	userRepo := repository.NewUserRepository(db)
	songRepo := repository.NewSongRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// 2. UseCases
	userUc := user.NewUserUseCase(userRepo)
	songUc := song.NewSongUseCase(songRepo)

	createPlaylistUc := playlist.NewCreatePlaylist(playlistRepo)
	fetchPlaylistsUc := playlist.NewFetchPlaylists(playlistRepo)
	addSongToPlaylistUc := playlist.NewAddSongToPlaylist(playlistRepo)

	// 3. Handlers
	return &AppHandlers{
		UserHandler: handlers.NewUserHandler(userUc),
		SongHandler: handlers.NewSongHandler(songUc),
		PlaylistHandler: handlers.NewPlaylistHandler(
			createPlaylistUc,
			fetchPlaylistsUc,
			addSongToPlaylistUc,
		),
	}
}
