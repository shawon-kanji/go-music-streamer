package handlers

import (
	"errors"
	"net/http"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
	"go-music-streamer/internal/framework"
	"go-music-streamer/internal/usecase/song"

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
			framework.NotFound(c)
			return
		}
		framework.InternalServerError(c, err)
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

	val, exists := c.Get("validatedRequest")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated request missing from context"))
		return
	}

	req, ok := val.(*dto.UpdateSongRequest)
	if !ok || req == nil {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated request has invalid type in context"))
		return
	}

	updatedSong, err := h.uc.UpdateSong(uriReq.ID, req)
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

type SongHandler struct {
	uc       song.SongUseCase
	tagQueue tagGeneratorQueue
}

type searchSongsUseCase interface {
	SearchSongs(req *dto.SearchSongsRequest) (*dto.SemanticSearchResponse, error)
}

type SearchSongsHandler struct {
	uc searchSongsUseCase
}

func NewSearchSongsHandler(uc searchSongsUseCase) *SearchSongsHandler {
	return &SearchSongsHandler{uc: uc}
}

func (h *SearchSongsHandler) Handle(c *gin.Context) {
	val, exists := c.Get("validatedRequest")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated request missing from context"))
		return
	}

	req, ok := val.(*dto.SearchSongsRequest)
	if !ok || req == nil {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated request has invalid type in context"))
		return
	}

	res, err := h.uc.SearchSongs(req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.Code == apperror.BadRequest {
			framework.SendError(c, http.StatusBadRequest, "Failed to search songs", err)
			return
		}

		framework.SendError(c, http.StatusInternalServerError, "Failed to search songs", err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Semantic search completed", res)
}

type tagGeneratorQueue interface {
	AddTask(task *entity.Song)
}

func NewSongHandler(uc song.SongUseCase, tagQueue tagGeneratorQueue) *SongHandler {
	return &SongHandler{uc: uc, tagQueue: tagQueue}
}

func (h *SongHandler) UploadSong(c *gin.Context) {
	req := dto.UploadSongRequest{}
	if err := c.ShouldBind(&req); err != nil {
		framework.SendError(c, http.StatusBadRequest, "Failed to bind request data", framework.FormatValidationError(err))
		return
	}

	entity, err := h.uc.UploadSong(&req, c)
	if err != nil {
		framework.SendError(c, http.StatusInternalServerError, "Failed to save uploaded file", err)
		return
	}

	if h.tagQueue != nil {
		h.tagQueue.AddTask(entity)
	}

	// For demonstration, we just return the file name. In a real application, you would save the file and create a song record.
	framework.SendSuccess(c, http.StatusOK, "File uploaded successfully", gin.H{"fileName": entity.URL})
}
