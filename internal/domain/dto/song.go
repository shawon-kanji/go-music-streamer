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
