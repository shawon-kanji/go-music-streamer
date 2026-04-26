package postgres

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string     `gorm:"unique;not null"`
	Email     string     `gorm:"unique;not null"`
	Password  string     `gorm:"not null"`
	UserRoles []UserRole `gorm:"foreignKey:UserID"`
	IsCreator bool       `gorm:"default:false"`
}

type UserRole struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	RoleID uint `gorm:"not null"`
	Role   Role `gorm:"foreignKey:RoleID"`
}

type AdminRole struct {
	gorm.Model
	AdminID uint `gorm:"not null"`
	RoleID  uint `gorm:"not null"`
	Role    Role `gorm:"foreignKey:RoleID"`
}

type Role struct {
	gorm.Model
	Name        string           `gorm:"unique;not null"`
	Permissions []RolePermission `gorm:"foreignKey:RoleID"`
}

// RolePermission ties a specific Action on a specific Resource to a Role
type RolePermission struct {
	gorm.Model
	RoleID     uint `gorm:"not null;uniqueIndex:idx_role_resource_action"`
	ResourceID uint `gorm:"not null;uniqueIndex:idx_role_resource_action"`
	ActionID   uint `gorm:"not null;uniqueIndex:idx_role_resource_action"`

	// Preload these to know exactly what the permission grants
	Resource Resource `gorm:"foreignKey:ResourceID"`
	Action   Action   `gorm:"foreignKey:ActionID"`
}

type Resource struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Action struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Admin struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Playlist struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Visibility  string `gorm:"not null; default:'public'"` // e.g., "public", "private"
	CreatedBy   uint   `gorm:"not null"`                   // UserID of the creator
	UserDetails User   `gorm:"foreignKey:CreatedBy"`
	Thumbnail   string // URL to playlist thumbnail image
	Songs       []Song `gorm:"many2many:playlist_songs;"`
}

type Song struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Artist    string `gorm:"not null"`
	Album     string
	Genre     string
	Duration  uint   // Duration in seconds
	Url       string `gorm:"not null"` // URL to the song file
	LikeCount uint   `gorm:"default:0"`
	Thumbnail string // URL to song thumbnail image
}

type UserLikedSong struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex:idx_user_liked_song"`
	SongID uint `gorm:"not null;uniqueIndex:idx_user_liked_song"`
}
