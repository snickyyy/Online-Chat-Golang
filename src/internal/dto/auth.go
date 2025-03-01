package dto

type RegisterRequest struct {
	Username 		string 	`json:"username" binding:"required,min=4,max=28"`
	Email       	string  `json:"email" binding:"required,email"`
	Password    	string  `json:"password" binding:"required"`
	ConfirmPassword string 	`json:"confirm_password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message" binding:"required"`
	Status  bool   `json:"status" binding:"required"`
}
