package models

import "time"

type PlaceOrderRequest struct {
	Name         string             `json:"name"`
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
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID           int       `json:"id" db:"id"`
	OwnerID      string    `json:"owner_id" db:"owner_id"`
	Total        float64   `json:"total" db:"total"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Status       string    `json:"status" db:"status"`
	Name         *string   `json:"name" db:"name"`
	Address      *string   `json:"address" db:"address"`
	DeliveryType *string   `json:"delivery_type" db:"delivery_type"`
	PaymentType  *string   `json:"payment_type" db:"payment_type"`
	ChangeFor    *float64  `json:"change_for" db:"change_for"`
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
