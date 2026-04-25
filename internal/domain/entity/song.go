package entity

type Song struct {
	ID        uint
	Title     string
	Artist    string
	Album     string
	Duration  int // Duration in seconds
	Url       string
	Genre     string
	LikeCount uint
}
