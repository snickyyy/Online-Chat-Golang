package dto

type ErrorResponse struct {
	Error string `json:"error" binding:"required"`
	Code  int    `json:"-"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
