package song

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (useCase *songUseCase) UploadSong(req *dto.UploadSongRequest, c *gin.Context) (*entity.Song, error) {
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(req.File.Filename))
	fileNameHash := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String(), ext)
	dstPath := filepath.Join("uploads", fileNameHash)

	if err := c.SaveUploadedFile(req.File, dstPath); err != nil {
		return nil, err
	}

	songEntity := &entity.Song{
		Title:    req.Title,
		Artist:   req.Artist,
		Album:    req.Album,
		Genre:    req.Genre,
		Duration: req.Duration,
		URL:      dstPath,
	}

	if entity, err := useCase.repo.CreateSong(songEntity); err != nil {
		return nil, err
	} else {
		return entity, nil
	}

}
