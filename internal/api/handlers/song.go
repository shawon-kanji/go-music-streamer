package handlers

import (
	"errors"
	"net/http"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/framework"

	"github.com/gin-gonic/gin"
)

type listSongsUseCase interface {
	ListSongs(page int, limit int) (*dto.PaginatedSongResponse, error)
}

type ListSongsHandler struct {
	uc listSongsUseCase
}

func NewListSongsHandler(uc listSongsUseCase) *ListSongsHandler {
	return &ListSongsHandler{uc: uc}
}

func (h *ListSongsHandler) Handle(c *gin.Context) {
	req := dto.ListSongsRequest{
		Page:  1,
		Limit: 10,
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	paginatedSongs, err := h.uc.ListSongs(req.Page, req.Limit)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Songs retrieved successfully", paginatedSongs)
}

type fetchSongUseCase interface {
	FetchSongByID(id uint) (*dto.SongResponse, error)
}

type FetchSongHandler struct {
	uc fetchSongUseCase
}

func NewFetchSongHandler(uc fetchSongUseCase) *FetchSongHandler {
	return &FetchSongHandler{uc: uc}
}

func (h *FetchSongHandler) Handle(c *gin.Context) {
	var req dto.FetchSongRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	songRes, err := h.uc.FetchSongByID(req.ID)
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

type updateSongUseCase interface {
	UpdateSong(id uint, req *dto.UpdateSongRequest) (*dto.SongResponse, error)
}

type UpdateSongHandler struct {
	uc updateSongUseCase
}

func NewUpdateSongHandler(uc updateSongUseCase) *UpdateSongHandler {
	return &UpdateSongHandler{uc: uc}
}

func (h *UpdateSongHandler) Handle(c *gin.Context) {
	var uriReq dto.FetchSongRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	var req dto.UpdateSongRequest
	val, exists := c.Get("validatedRequest")
	if exists {
		req = *val.(*dto.UpdateSongRequest)
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}

	updatedSong, err := h.uc.UpdateSong(uriReq.ID, &req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.Code == apperror.NotFound {
			framework.SendError(c, http.StatusNotFound, err.Error(), err)
			return
		}
		framework.SendError(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Song updated successfully", updatedSong)
}
