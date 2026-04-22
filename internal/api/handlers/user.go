package handlers

import (
	"net/http"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/usecase/user" // Assuming package is renamed to user

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCase user.UserUseCase
}

func NewUserHandler(useCase user.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var createReq dto.CreateUserRequest

	// Bind incoming JSON to the Request DTO
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Execute UseCase (which now inherently returns the sanitized Response DTO)
	userResp, err := h.useCase.CreateUser(&createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    userResp,
	})
}

// Login handler parsing credentials and returning the Token securely via the UseCase layer
func (h *UserHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	token, userResp, err := h.useCase.Login(&loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    *userResp,
	})
}
