package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-music-streamer/internal/domain/apperror"
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

type updatePlaylistUseCase interface {
	UpdatePlaylist(playlistID uint, userID uint, req *dto.UpdatePlaylistRequest) (*dto.PlaylistResponse, error)
}

type UpdatePlaylistHandler struct {
	uc updatePlaylistUseCase
}

func NewUpdatePlaylistHandler(uc updatePlaylistUseCase) *UpdatePlaylistHandler {
	return &UpdatePlaylistHandler{uc: uc}
}

func (h *UpdatePlaylistHandler) Handle(c *gin.Context) {
	playlistIDStr := c.Param("id")
	playlistID, err := strconv.ParseUint(playlistIDStr, 10, 32)
	if err != nil {
		framework.BadRequest(c, fmt.Errorf("invalid playlist ID: %v", err))
		return
	}

	var req dto.UpdatePlaylistRequest
	val, exists := c.Get("validatedRequest")
	if exists {
		req = *val.(*dto.UpdatePlaylistRequest)
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	res, err := h.uc.UpdatePlaylist(uint(playlistID), userID, &req)
	if err != nil {
		handlePlaylistEditError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Playlist updated successfully", res)
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
	AddSong(playlistID uint, userID uint, req *dto.AddSongToPlaylistRequest) error
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

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	err = h.uc.AddSong(uint(playlistID), userID, &req)
	if err != nil {
		handlePlaylistEditError(c, err)
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

func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "identity missing from secure context"))
		return 0, false
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "invalid identity format in secure context"))
		return 0, false
	}

	return userID, true
}

func handlePlaylistEditError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case apperror.NotFound:
			framework.SendError(c, http.StatusNotFound, appErr.Message, err)
			return
		case apperror.Unauthorized:
			framework.Forbidden(c, err)
			return
		case apperror.BadRequest, apperror.DataValidationError:
			framework.BadRequest(c, err)
			return
		}
	}

	framework.InternalServerError(c, err)
}
