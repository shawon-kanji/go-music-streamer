package bootstrap

import (
	"go-music-streamer/internal/api/handlers"
	"go-music-streamer/internal/repository"
	"go-music-streamer/internal/schedulers"
	"go-music-streamer/internal/usecase/playlist"
	"go-music-streamer/internal/usecase/song"
	"go-music-streamer/internal/usecase/user"

	"gorm.io/gorm"
)

// AppHandlers struct acts as a container for all HTTP handlers
type AppHandlers struct {
	UserSignupHandler  *handlers.UserSignupHandler
	UserLoginHandler   *handlers.UserLoginHandler
	UserProfileHandler *handlers.UserProfileHandler

	ListSongsHandler  *handlers.ListSongsHandler
	FetchSongHandler  *handlers.FetchSongHandler
	UpdateSongHandler *handlers.UpdateSongHandler
	SongHandler       *handlers.SongHandler

	CreatePlaylistHandler    *handlers.CreatePlaylistHandler
	UpdatePlaylistHandler    *handlers.UpdatePlaylistHandler
	FetchPlaylistsHandler    *handlers.FetchPlaylistsHandler
	FetchPlaylistHandler     *handlers.FetchPlaylistHandler
	AddSongToPlaylistHandler *handlers.AddSongToPlaylistHandler
}

// InitHandlers initializes and injects all repositories, use cases, and handlers.
func InitHandlers(db *gorm.DB) *AppHandlers {
	// 1. Repositories
	userRepo := repository.NewUserRepository(db)
	songRepo := repository.NewSongRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// 2. UseCases (Unified)
	userUc := user.NewUserUseCase(userRepo)
	songUc := song.NewSongUseCase(songRepo)
	playlistUc := playlist.NewUseCase(playlistRepo)

	// 3. Schedulers
	scheduler := &schedulers.Scheduler{}
	schedulers.InitiateSchedulers(scheduler)

	// 4. Handlers (Granular, decoupled)
	return &AppHandlers{
		UserSignupHandler:  handlers.NewUserSignupHandler(userUc),
		UserLoginHandler:   handlers.NewUserLoginHandler(userUc),
		UserProfileHandler: handlers.NewUserProfileHandler(),

		ListSongsHandler:  handlers.NewListSongsHandler(songUc),
		FetchSongHandler:  handlers.NewFetchSongHandler(songUc),
		UpdateSongHandler: handlers.NewUpdateSongHandler(songUc),
		SongHandler:       handlers.NewSongHandler(songUc, scheduler.TagGenerator),

		CreatePlaylistHandler:    handlers.NewCreatePlaylistHandler(playlistUc),
		UpdatePlaylistHandler:    handlers.NewUpdatePlaylistHandler(playlistUc),
		FetchPlaylistsHandler:    handlers.NewFetchPlaylistsHandler(playlistUc),
		FetchPlaylistHandler:     handlers.NewFetchPlaylistHandler(playlistUc),
		AddSongToPlaylistHandler: handlers.NewAddSongToPlaylistHandler(playlistUc),
	}
}
