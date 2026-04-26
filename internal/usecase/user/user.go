package user

import (
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/repository"
)

type UserUseCase interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (string, *dto.UserResponse, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}
