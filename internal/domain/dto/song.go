package dto

type ListSongsRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type SongResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Duration  uint   `json:"duration"` // Duration in seconds
	URL       string `json:"url"`
	LikeCount uint   `json:"like_count"`
	Genre     string `json:"genre"`
	Thumbnail string `json:"thumbnail"`
}

type PaginatedSongResponse struct {
	MediaList  []*SongResponse `json:"media_list"`
	HasMore    bool            `json:"has_more"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalCount int64           `json:"total_count"`
}

type FetchSongRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdateSongRequest struct {
	Title     *string `json:"title" binding:"omitempty"`
	Artist    *string `json:"artist" binding:"omitempty"`
	Album     *string `json:"album" binding:"omitempty"`
	Duration  *uint   `json:"duration" binding:"omitempty"`
	URL       *string `json:"url" binding:"omitempty"`
	Genre     *string `json:"genre" binding:"omitempty"`
	Thumbnail *string `json:"thumbnail" binding:"omitempty"`
}
