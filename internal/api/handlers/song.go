package handlers

import (
	"errors"
	"net/http"

	"go-music-streamer/internal/domain/apperror"
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

func (h *SongHandler) FetchSong(c *gin.Context) {
	var req dto.FetchSongRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	songRes, err := h.useCase.FetchSongByID(req.ID)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.Code == apperror.NotFound {
			framework.SendError(c, http.StatusNotFound, err.Error(), err)
			return
		}
		framework.SendError(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Song retrieved successfully", songRes)
}
