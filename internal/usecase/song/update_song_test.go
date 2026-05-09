package song

import (
	"errors"
	"testing"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
)

type fakeSongRepo struct {
	createSongFn  func(song *entity.Song) (*entity.Song, error)
	getSongFn     func(id uint) (*entity.Song, error)
	updateSongFn  func(song *entity.Song) (*entity.Song, error)
	updateCalled  bool
	lastUpdateArg *entity.Song
}

func (f *fakeSongRepo) CreateSong(song *entity.Song) (*entity.Song, error) {
	if f.createSongFn != nil {
		return f.createSongFn(song)
	}
	return song, nil
}

func (f *fakeSongRepo) ListSongs(page int, limit int) ([]*entity.Song, int64, error) {
	return nil, 0, nil
}

func (f *fakeSongRepo) GetSongByID(id uint) (*entity.Song, error) {
	if f.getSongFn != nil {
		return f.getSongFn(id)
	}
	return nil, nil
}

func (f *fakeSongRepo) UpdateSong(song *entity.Song) (*entity.Song, error) {
	f.updateCalled = true
	f.lastUpdateArg = song
	if f.updateSongFn != nil {
		return f.updateSongFn(song)
	}
	return song, nil
}

func TestUpdateSong_Success(t *testing.T) {
	repo := &fakeSongRepo{}
	repo.getSongFn = func(id uint) (*entity.Song, error) {
		return &entity.Song{
			ID:        id,
			Title:     "old",
			Artist:    "old artist",
			Album:     "old album",
			Genre:     "old genre",
			Duration:  100,
			URL:       "old-url",
			Thumbnail: "old-thumb",
		}, nil
	}

	uc := &songUseCase{repo: repo}
	newTitle := "new"
	newURL := "new-url"
	req := &dto.UpdateSongRequest{
		Title: &newTitle,
		URL:   &newURL,
	}

	res, err := uc.UpdateSong(7, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.updateCalled {
		t.Fatalf("expected UpdateSong repository call")
	}
	if repo.lastUpdateArg == nil {
		t.Fatalf("expected update argument to be captured")
	}
	if repo.lastUpdateArg.Title != "new" {
		t.Fatalf("expected title to be updated, got %q", repo.lastUpdateArg.Title)
	}
	if repo.lastUpdateArg.URL != "new-url" {
		t.Fatalf("expected url to be updated, got %q", repo.lastUpdateArg.URL)
	}
	if repo.lastUpdateArg.Artist != "old artist" {
		t.Fatalf("expected artist unchanged, got %q", repo.lastUpdateArg.Artist)
	}
	if res == nil || res.Title != "new" || res.URL != "new-url" {
		t.Fatalf("expected updated response values")
	}
}

func TestUpdateSong_GetSongError(t *testing.T) {
	repo := &fakeSongRepo{}
	expectedErr := errors.New("load failed")
	repo.getSongFn = func(id uint) (*entity.Song, error) {
		return nil, expectedErr
	}

	uc := &songUseCase{repo: repo}
	_, err := uc.UpdateSong(7, &dto.UpdateSongRequest{})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if repo.updateCalled {
		t.Fatalf("did not expect update call when load fails")
	}
}
