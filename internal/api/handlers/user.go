package handlers

import (
	"net/http"

	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/framework"
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
		errs := framework.FormatValidationError(err)
		framework.BadRequest(c, errs)
		return
	}

	// Execute UseCase (which now inherently returns the sanitized Response DTO)
	userResp, err := h.useCase.CreateUser(&createReq)
	if err != nil {
		framework.BadRequest(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusCreated, "User created successfully", userResp)
}

// Login handler parsing credentials and returning the Token securely via the UseCase layer
func (h *UserHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		errs := framework.FormatValidationError(err)
		framework.BadRequest(c, errs)
		return
	}

	token, userResp, err := h.useCase.Login(&loginReq)
	if err != nil {
		framework.Unauthorized(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Login successful", dto.LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    *userResp,
	})
}
