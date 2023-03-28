package dto

type SignavaultError struct {
	Message string `json:"message" binding:"required"`
	Error   string `json:"error" binding:"required"`
}
