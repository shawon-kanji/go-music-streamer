package dto

// CreateUserRequest is purely for incoming data
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse is purely for outgoing data (NO password field exists here)
type UserResponse struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	UserRoles []string `json:"user_roles,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
