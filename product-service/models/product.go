package models

import "time"

type Product struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Price     float64    `db:"price"`
	Quantity  int        `db:"quantity"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
