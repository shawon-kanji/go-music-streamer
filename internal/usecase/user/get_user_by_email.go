package user

import (
	"go-music-streamer/internal/domain/dto"
)

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
