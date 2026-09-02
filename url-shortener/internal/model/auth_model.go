package model

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" db:"email"`
	Password string `json:"password" binding:"required,min=8" db:"password"`
}
