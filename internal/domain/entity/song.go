package entity

type Song struct {
	ID        uint
	Title     string
	Artist    string
	Album     string
	Duration  uint // Duration in seconds
	URL       string
	Genre     string
	LikeCount uint
	Thumbnail string
}
