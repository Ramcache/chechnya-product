package models

import "time"

type Product struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	Price        float64   `db:"price" json:"price"`
	Availability bool      `db:"availability" json:"availability"`
	CategoryID   int       `db:"category_id" json:"category_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type ProductResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Availability bool    `json:"availability"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Rating       float64 `json:"rating"`
}
