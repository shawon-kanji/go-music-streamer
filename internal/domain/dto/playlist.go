package dto

type Playlist struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Visibility  string         `json:"visibility"`
	CreatedBy   uint           `json:"created_by"`
	UserDetails UserResponse   `json:"user_details"`
	Songs       []SongResponse `json:"songs"`
}

type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Visibility  string `json:"visibility" binding:"required,oneof=public private"`
	Thumbnail   string `json:"thumbnail"`
	CreatedBy   uint   `json:"created_by"`
}

type UpdatePlaylistRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Visibility  string `json:"visibility" binding:"omitempty,oneof=public private"`
	Thumbnail   string `json:"thumbnail" binding:"omitempty"`
}

type AddSongToPlaylistRequest struct {
	SongID uint `json:"song_id" binding:"required"`
}

type RemoveSongFromPlaylistRequest struct {
	SongID uint `json:"song_id" binding:"required"`
}

type PlaylistUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PlaylistResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Visibility  string         `json:"visibility"`
	CreatedBy   uint           `json:"created_by"`
	UserDetails PlaylistUser   `json:"user_details"`
	Thumbnail   string         `json:"thumbnail"`
	Songs       []SongResponse `json:"songs"`
}

type PaginatedPlaylistResponse struct {
	PlaylistList []*PlaylistResponse `json:"playlist_list"`
	HasMore      bool                `json:"has_more"`
	Page         int                 `json:"page"`
	Limit        int                 `json:"limit"`
	TotalCount   int64               `json:"total_count"`
}
