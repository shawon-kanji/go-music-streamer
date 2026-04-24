package postgres

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string     `gorm:"unique;not null"`
	Email     string     `gorm:"unique;not null"`
	Password  string     `gorm:"not null"`
	UserRoles []UserRole `gorm:"foreignKey:UserID"`
}

type UserRole struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	RoleID uint `gorm:"not null"`
	Role   Role `gorm:"foreignKey:RoleID"`
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
