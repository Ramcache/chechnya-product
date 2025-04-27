package models

import (
	"database/sql"
	"time"
)

type Product struct {
	ID           int            `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Description  string         `db:"description" json:"description"`
	Price        float64        `db:"price" json:"price"`
	Availability bool           `db:"availability" json:"availability"`
	CategoryID   sql.NullInt64  `db:"category_id" json:"category_id"`
	Url          sql.NullString `db:"url" json:"url"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
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
	Url          string  `json:"url"`
}

type ProductInput struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Availability bool    `json:"availability"`
	CategoryID   *int    `json:"category_id"`
	Url          string  `json:"url"`
}

type ProductPatchInput struct {
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	Availability *bool    `json:"availability,omitempty"`
	CategoryID   *int     `json:"category_id,omitempty"`
	Url          *string  `json:"url,omitempty"`
}
