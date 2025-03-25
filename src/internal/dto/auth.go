package dto

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,username"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,password"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message" binding:"required"`
	Status  bool   `json:"status" binding:"required"`
}
