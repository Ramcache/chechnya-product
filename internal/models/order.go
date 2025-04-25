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

type OrderItem struct {
	OrderID   int `db:"order_id" json:"order_id"`
	ProductID int `db:"product_id" json:"product_id"`
	Quantity  int `db:"quantity" json:"quantity"`
}

type OrderWithItems struct {
	ID        int             `json:"id"`
	OwnerID   string          `json:"owner_id"`
	Total     float64         `json:"total"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"created_at"`
	Items     []OrderItemFull `json:"items"`
}

type OrderItemFull struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
