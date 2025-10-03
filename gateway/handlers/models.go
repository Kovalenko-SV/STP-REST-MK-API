package handlers

import "time"

type User struct {
	ID        string     `db:"id"`
	Login     string     `db:"login"`
	Password  string     `db:"password"`
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Product struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Price     float64    `db:"price" json:"price"`
	Quantity  int        `db:"quantity" json:"quantity"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID      string `json:"id"`
	Login   string `json:"login"`
	IsAdmin bool   `json:"is_admin"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}
