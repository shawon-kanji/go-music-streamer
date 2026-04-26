package handlers

import (
	"net/http"

	"go-music-streamer/internal/domain/apperror"
	"go-music-streamer/internal/domain/dto"
	"go-music-streamer/internal/framework"

	"github.com/gin-gonic/gin"
)

type createUserUseCase interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
}

type UserSignupHandler struct {
	uc createUserUseCase
}

func NewUserSignupHandler(uc createUserUseCase) *UserSignupHandler {
	return &UserSignupHandler{uc: uc}
}

func (h *UserSignupHandler) Handle(c *gin.Context) {
	val, exists := c.Get("validatedRequest")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated data not found in context"))
		return
	}

	createReq := val.(*dto.CreateUserRequest)

	userResp, err := h.uc.CreateUser(createReq)
	if err != nil {
		framework.BadRequest(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusCreated, "User created successfully", userResp)
}

type loginUserUseCase interface {
	Login(req *dto.LoginRequest) (string, *dto.UserResponse, error)
}

type UserLoginHandler struct {
	uc loginUserUseCase
}

func NewUserLoginHandler(uc loginUserUseCase) *UserLoginHandler {
	return &UserLoginHandler{uc: uc}
}

func (h *UserLoginHandler) Handle(c *gin.Context) {
	val, exists := c.Get("validatedRequest")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "validated data not found in context"))
		return
	}

	loginReq := val.(*dto.LoginRequest)

	token, userResp, err := h.uc.Login(loginReq)
	if err != nil {
		framework.Unauthorized(c, err)
		return
	}

	framework.SendSuccess(c, http.StatusOK, "Login successful", dto.LoginResponse{
		Token: token,
		User:  *userResp,
	})
}

// UserProfileHandler doesn't even need a usecase right now since it just returns the userID,
// but we'll create the structure for consistency.
type UserProfileHandler struct {
}

func NewUserProfileHandler() *UserProfileHandler {
	return &UserProfileHandler{}
}

func (h *UserProfileHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		framework.InternalServerError(c, apperror.New(apperror.InternalError, "Identity missing from secure context"))
		return
	}

	framework.SendSuccess(c, http.StatusOK, "User profile retrieved", gin.H{
		"id": userID,
	})
}
