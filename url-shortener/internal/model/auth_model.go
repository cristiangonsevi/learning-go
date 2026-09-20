package model

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" db:"email"`
	Password string `json:"password" binding:"required,min=8" db:"password"`
}

type LoginResponse struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
