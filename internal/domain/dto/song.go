package dto

type SongResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Duration  int    `json:"duration"` // Duration in seconds
	URL       string `json:"url"`
	LikeCount uint   `json:"like_count"`
	Genre     string `json:"genre"`
}
