package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/framework"
	"go-music-streamer/internal/usecase/playlist"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	createPlaylist    playlist.CreatePlaylist
	fetchPlaylists    playlist.FetchPlaylists
	addSongToPlaylist playlist.AddSongToPlaylist
}

func NewPlaylistHandler(
	cp playlist.CreatePlaylist,
	fp playlist.FetchPlaylists,
	as playlist.AddSongToPlaylist,
) *PlaylistHandler {
	return &PlaylistHandler{
		createPlaylist:    cp,
		fetchPlaylists:    fp,
		addSongToPlaylist: as,
	}
}

// CreatePlaylist generates a new playlist for the logged-in user
func (h *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	var req dto.CreatePlaylistRequest

	val, exists := c.Get("validatedRequest")
	if exists {
		req = *val.(*dto.CreatePlaylistRequest)
	} else if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, framework.FormatValidationError(err))
		return
	}
	req.CreatedBy = c.MustGet("userID").(uint) // Override with authenticated user ID from context

	res, err := h.createPlaylist.CreatePlaylist(&req)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusCreated, "Playlist created successfully", res)
}

// FetchPlaylists retrieves paginated playlists
func (h *PlaylistHandler) FetchPlaylists(c *gin.Context) {
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

	res, err := h.fetchPlaylists.FetchPlaylists(req.Page, req.Limit)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Playlists fetched successfully", res)
}

// AddSong adds a single song to a specified playlist via path
func (h *PlaylistHandler) AddSong(c *gin.Context) {
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

	err = h.addSongToPlaylist.AddSong(uint(playlistID), &req)
	if err != nil {
		framework.InternalServerError(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Song added to playlist successfully", nil)
}
