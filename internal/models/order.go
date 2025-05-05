package models

import (
	"time"
)

type PlaceOrderRequest struct {
	Name         *string            `json:"name"`
	Address      *string            `json:"address"`
	Items        []OrderItemRequest `json:"items"`
	Total        float64            `json:"total"`
	PaymentType  string             `json:"paymentType"`
	Status       string             `json:"status"`
	DeliveryType string             `json:"deliveryType"`
	CreatedAt    int64              `json:"createdAt"`
	DeliveryText string             `json:"deliveryText"`
	DeliveryFee  float64            `json:"deliveryFee"`
	ChangeFor    *float64           `json:"changeFor"`
}

type OrderItemRequest struct {
	ID       int     `json:"id"`
	Name     *string `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID                int       `db:"id"`
	OwnerID           string    `db:"owner_id"`
	Total             float64   `db:"total"`
	CreatedAt         time.Time `db:"created_at"`
	Status            string    `db:"status"`
	Name              *string   `db:"name"`
	Address           *string   `db:"address"`
	DeliveryType      string    `db:"delivery_type"`
	PaymentType       string    `db:"payment_type"`
	ChangeFor         *float64  `db:"change_for"`
	DeliveryFee       *float64  `db:"delivery_fee"`
	DeliveryText      *string   `db:"delivery_text"`
	FrontendCreatedAt *int64    `db:"frontend_created_at"`
}

type OrderStatusRequest struct {
	Status string `json:"status" example:"в пути"`
}

type OrderItem struct {
	OrderID   int     `db:"order_id" json:"order_id"`
	Name      *string `json:"name" db:"name"`
	ProductID int     `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
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
