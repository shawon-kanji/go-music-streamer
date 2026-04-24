package repository

import (
	"go-music-streamer/internal/database/postgres"
	"go-music-streamer/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entity.User) error {
	// Map Domain Entity to Database Schema Model
	dbUser := &postgres.User{
		Username: user.Name,
		Email:    user.Email,
		Password: user.Password, // Ideally, hashed before reaching the repository
	}

	if err := r.db.Create(dbUser).Error; err != nil {
		return err
	}

	// Update the domain entity with the generated ID
	user.ID = dbUser.ID
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	var dbUser postgres.User

	// Fetch from DB using GORM
	if err := r.db.Where(map[string]interface{}{"email": email}).First(&dbUser).Error; err != nil {
		return nil, err
	}

	// Map Database Schema Model back to Domain Entity
	return &entity.User{
		ID:       dbUser.ID,
		Name:     dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}, nil
}
