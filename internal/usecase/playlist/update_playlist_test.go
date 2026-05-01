package playlist
package playlist

import (
	"errors"
	"testing"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"

	"gorm.io/gorm"
)

type fakePlaylistRepo struct {
	getPlaylistFn  func(id uint) (*entity.Playlist, error)
	updateFn       func(playlist *entity.Playlist) (entity.Playlist, error)
	updateCalled   bool
	lastUpdateArg  *entity.Playlist
}

func (f *fakePlaylistRepo) CreatePlaylist(playlist *entity.Playlist) (entity.Playlist, error) {
	return entity.Playlist{}, nil
}

func (f *fakePlaylistRepo) GetPlaylistByID(id uint) (*entity.Playlist, error) {
	if f.getPlaylistFn != nil {
		return f.getPlaylistFn(id)
	}
	return nil, nil
}

func (f *fakePlaylistRepo) UpdatePlaylist(playlist *entity.Playlist) (entity.Playlist, error) {
	f.updateCalled = true
	f.lastUpdateArg = playlist
	if f.updateFn != nil {
		return f.updateFn(playlist)
	}
	return *playlist, nil
}

func (f *fakePlaylistRepo) DeletePlaylist(id uint) error {
	return nil
}

func (f *fakePlaylistRepo) ListPlaylists(page int, limit int, options dto.PlaylistQueryOptions) ([]*entity.Playlist, int64, error) {
	return nil, 0, nil
}

func (f *fakePlaylistRepo) AddSongToPlaylist(playlistID uint, songID uint) error {
	return nil
}

func TestUpdatePlaylist_Success(t *testing.T) {
	repo := &fakePlaylistRepo{}
	repo.getPlaylistFn = func(id uint) (*entity.Playlist, error) {
		return &entity.Playlist{
			ID:          id,
			Name:        "old name",
			Description: "old desc",
			Visibility:  "private",
			CreatedBy:   42,
			Thumbnail:   "old-thumb",
		}, nil
	}

	uc := &useCase{playlistRepo: repo}
	req := &dto.UpdatePlaylistRequest{
		Name:       "new name",
		Visibility: "public",
	}

	res, err := uc.UpdatePlaylist(10, 42, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.updateCalled {
		t.Fatalf("expected UpdatePlaylist repository call")
	}
	if repo.lastUpdateArg == nil {
		t.Fatalf("expected update argument to be captured")
	}
	if repo.lastUpdateArg.Name != "new name" {
		t.Fatalf("expected name to be updated, got %q", repo.lastUpdateArg.Name)
	}
	if repo.lastUpdateArg.Visibility != "public" {
		t.Fatalf("expected visibility to be updated, got %q", repo.lastUpdateArg.Visibility)
	}
	if repo.lastUpdateArg.Description != "old desc" {
		t.Fatalf("expected description unchanged, got %q", repo.lastUpdateArg.Description)
	}
	if res == nil || res.Name != "new name" {
		t.Fatalf("expected response name to be updated")
	}
}

func TestUpdatePlaylist_Unauthorized(t *testing.T) {
	repo := &fakePlaylistRepo{}
	repo.getPlaylistFn = func(id uint) (*entity.Playlist, error) {
		return &entity.Playlist{ID: id, CreatedBy: 99}, nil
	}

	uc := &useCase{playlistRepo: repo}
	_, err := uc.UpdatePlaylist(10, 42, &dto.UpdatePlaylistRequest{Name: "x"})
	if err == nil {
		t.Fatalf("expected unauthorized error")
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != apperror.Unauthorized {
		t.Fatalf("expected code %s, got %s", apperror.Unauthorized, appErr.Code)
	}
	if repo.updateCalled {
		t.Fatalf("did not expect update call when unauthorized")
	}
}

func TestUpdatePlaylist_NotFound(t *testing.T) {
	repo := &fakePlaylistRepo{}
	repo.getPlaylistFn = func(id uint) (*entity.Playlist, error) {
		return nil, gorm.ErrRecordNotFound
	}

	uc := &useCase{playlistRepo: repo}
	_, err := uc.UpdatePlaylist(10, 42, &dto.UpdatePlaylistRequest{Name: "x"})
	if err == nil {
		t.Fatalf("expected not found error")
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != apperror.NotFound {
		t.Fatalf("expected code %s, got %s", apperror.NotFound, appErr.Code)
	}
}
