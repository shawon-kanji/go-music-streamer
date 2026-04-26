package entity

type Playlist struct {
	ID          uint
	Name        string
	Description string
	Visibility  string
	CreatedBy   uint
	UserDetails User
	Thumbnail   string
	Songs       []Song
}
