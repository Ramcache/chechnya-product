package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

type OrderStatusRequest struct {
	Status string `json:"status" example:"в пути"`
}
