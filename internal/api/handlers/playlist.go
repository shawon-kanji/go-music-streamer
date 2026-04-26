package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/framework"

	"github.com/gin-gonic/gin"
)

type createPlaylistUseCase interface {
	CreatePlaylist(req *dto.CreatePlaylistRequest) (*dto.PlaylistResponse, error)
}

type CreatePlaylistHandler struct {
	uc createPlaylistUseCase
}

func NewCreatePlaylistHandler(uc createPlaylistUseCase) *CreatePlaylistHandler {
	return &CreatePlaylistHandler{uc: uc}
}

func (h *CreatePlaylistHandler) Handle(c *gin.Context) {
	var req dto.CreatePlaylistRequest

	val, exists := c.Get("validatedRequest")
	if exists {
		req = *val.(*dto.CreatePlaylistRequest)
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}
	req.CreatedBy = c.MustGet("userID").(uint) // Override with authenticated user ID from context

	res, err := h.uc.CreatePlaylist(&req)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusCreated, "Playlist created successfully", res)
}

type fetchPlaylistsUseCase interface {
	FetchPlaylists(page int, limit int, options dto.PlaylistQueryOptions) (*dto.PaginatedPlaylistResponse, error)
}

type FetchPlaylistsHandler struct {
	uc fetchPlaylistsUseCase
}

func NewFetchPlaylistsHandler(uc fetchPlaylistsUseCase) *FetchPlaylistsHandler {
	return &FetchPlaylistsHandler{uc: uc}
}

func (h *FetchPlaylistsHandler) Handle(c *gin.Context) {
	req := struct {
		Page  int `form:"page" binding:"omitempty,min=1"`
		Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
	}{
		Page:  1,
		Limit: 10,
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	res, err := h.uc.FetchPlaylists(req.Page, req.Limit, dto.PlaylistQueryOptions{
		UserID: c.MustGet("userID").(uint),
	})
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Playlists fetched successfully", res)
}

type addSongToPlaylistUseCase interface {
	AddSong(playlistID uint, req *dto.AddSongToPlaylistRequest) error
}

type AddSongToPlaylistHandler struct {
	uc addSongToPlaylistUseCase
}

func NewAddSongToPlaylistHandler(uc addSongToPlaylistUseCase) *AddSongToPlaylistHandler {
	return &AddSongToPlaylistHandler{uc: uc}
}

func (h *AddSongToPlaylistHandler) Handle(c *gin.Context) {
	playlistIDStr := c.Param("id")
	playlistID, err := strconv.ParseUint(playlistIDStr, 10, 32)
	if err != nil {
		framework.BadRequest(c, fmt.Errorf("invalid playlist ID: %v", err))
		return
	}

	var req dto.AddSongToPlaylistRequest
	val, exists := c.Get("validatedRequest")
	if exists {
		req = *val.(*dto.AddSongToPlaylistRequest)
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	err = h.uc.AddSong(uint(playlistID), &req)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Song added to playlist successfully", nil)
}

type fetchPlaylistUseCase interface {
	FetchPlaylist(id uint) (*dto.PlaylistResponse, error)
}

type FetchPlaylistHandler struct {
	uc fetchPlaylistUseCase
}

func NewFetchPlaylistHandler(uc fetchPlaylistUseCase) *FetchPlaylistHandler {
	return &FetchPlaylistHandler{uc: uc}
}

func (h *FetchPlaylistHandler) Handle(c *gin.Context) {
	playlistIDStr := c.Param("id")
	playlistID, err := strconv.ParseUint(playlistIDStr, 10, 32)
	if err != nil {
		framework.BadRequest(c, fmt.Errorf("invalid playlist ID"))
		return
	}

	res, err := h.uc.FetchPlaylist(uint(playlistID))
	if err != nil {
		framework.SendError(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Playlist fetched successfully", res)
}
