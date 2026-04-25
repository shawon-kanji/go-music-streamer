package handlers

import (
	"net/http"

	"go-music-streamer/internal/domain/apperror"
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
	// Retrieve the validated struct from the context set by the middleware
	val, exists := c.Get("validatedRequest")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated data not found in context"))
		return
	}

	createReq := val.(*dto.CreateUserRequest)

	// Execute UseCase (which now inherently returns the sanitized Response DTO)
	userResp, err := h.useCase.CreateUser(createReq)
	if err != nil {
		framework.BadRequest(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusCreated, "User created successfully", userResp)
}

// Login handler parsing credentials and returning the Token securely via the UseCase layer
func (h *UserHandler) Login(c *gin.Context) {
	val, exists := c.Get("validatedRequest")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated data not found in context"))
		return
	}

	loginReq := val.(*dto.LoginRequest)

	token, userResp, err := h.useCase.Login(loginReq)
	if err != nil {
		framework.Unauthorized(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Login successful", dto.LoginResponse{
		Token: token,
		User:  *userResp,
	})
}

// Profile returns the currently authenticated user's profile details.
func (h *UserHandler) Profile(c *gin.Context) {
	// The Authenticate middleware guarantees "userID" exists in the context
	userID, exists := c.Get("userID")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "Identity missing from secure context"))
		return
	}

	// For now, return a placeholder profile.
	// You can later map this to h.useCase.GetUserByID(userID)
	framework.SendSuccess(c, http.StatusOK, "User profile retrieved", gin.H{
		"id": userID,
	})
}
