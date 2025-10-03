package models

import "time"

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
