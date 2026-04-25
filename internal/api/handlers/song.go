package handlers

import (
	"go-music-streamer/internal/framework"
	"go-music-streamer/internal/usecase/song"
	"net/http"

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
	songs, err := h.useCase.ListSongs()
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Songs retrieved successfully", songs)
}
