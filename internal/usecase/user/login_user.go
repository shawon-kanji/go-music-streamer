package user

import (
	"errors"
	"time"

	"go-music-streamer/internal/config"
	"go-music-streamer/internal/domain/dto"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Login orchestrates fetching a user, computing password match, and creating a short-lived JWT token
func (useCase *userUseCase) Login(req *dto.LoginRequest) (string, *dto.UserResponse, error) {
	// Verify user exists via repo
	userEntity, err := useCase.repo.GetUserByEmail(req.Email)
	if err != nil { // Could be not found or other errors
		return "", nil, errors.New("invalid email or password")
	}

	// Compare bcrypt password
	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Fetch singleton config directly
	appConfig := config.GetConfig()

	// Generate Short-lived JWT properly sized with best practices (iss, aud, etc.)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userEntity.ID,
		"email": userEntity.Email,
		"aud":   "go-music-streamer-api",
		"iss":   "go-music-streamer",
		"exp":   time.Now().Add(15 * time.Minute).Unix(), // 15 mins expiration (Short lived)
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(), // Not before
	})

	tokenString, err := token.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return "", nil, errors.New("could not generate token")
	}

	// Return strictly safe UserResponse and Token
	respDTO := &dto.UserResponse{
		ID:    userEntity.ID,
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}

	return tokenString, respDTO, nil
}
