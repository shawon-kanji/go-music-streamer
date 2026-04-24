package user

import (
	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/domain/entity"
	"go-music-streamer/internal/repository"

	"golang.org/x/crypto/bcrypt"
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

// CreateUser handles business logic mapping the incoming DTO to a domain entity
// and delegating persistence to the repository layer.
func (useCase *userUseCase) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {

	//check if the user with the same email already exists
	existingUser, err := useCase.repo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, apperror.Newf("USER_ALREADY_EXISTS", "user with email %s already exists", req.Email)
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Map DTO -> Entity
	userEntity := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = useCase.repo.CreateUser(userEntity)
	if err != nil {
		return nil, err
	}

	// Map Entity -> Response DTO (Password omitted by design)
	return &dto.UserResponse{
		ID:    userEntity.ID,
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}, nil
}

// GetUserByEmail fetches a user via the repository and maps the returning entity back to a DTO
func (useCase *userUseCase) GetUserByEmail(email string) (*dto.UserResponse, error) {
	userEntity, err := useCase.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Map Entity -> Response DTO
	return &dto.UserResponse{
		ID:    userEntity.ID,
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}, nil
}
