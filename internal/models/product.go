package models

import (
	"time"
)

type Product struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Category    string    `json:"category" db:"category"`
	Stock       int       `db:"stock"`
	CreatedAt   time.Time `db:"created_at"`
}
