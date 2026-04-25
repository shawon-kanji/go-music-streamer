package handlers

import (
	"net/http"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/framework"
	"go-music-streamer/internal/usecase/song"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	useCase song.SongUseCase
}

func NewSongHandler(userCase song.SongUseCase) *SongHandler {
	return &SongHandler{
		useCase: userCase,
	}
}

func (h *SongHandler) ListSongs(c *gin.Context) {
	// Initialize with defaults before binding
	req := dto.ListSongsRequest{
		Page:  1,
		Limit: 10,
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	paginatedSongs, err := h.useCase.ListSongs(req.Page, req.Limit)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Songs retrieved successfully", paginatedSongs)
}
