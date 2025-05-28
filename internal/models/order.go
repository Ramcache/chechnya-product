package models

import (
	"encoding/json"
	"time"
)

// PlaceOrderRequest запрос от клиента для оформления заказа
type PlaceOrderRequest struct {
	Name         *string     `json:"name"`
	Address      *string     `json:"address"`
	Items        []OrderItem `json:"items"`
	PaymentType  string      `json:"payment_type"`
	Status       string      `json:"status"`
	DeliveryType string      `json:"delivery_type"`
	CreatedAt    int64       `json:"created_at"`
	DeliveryText string      `json:"delivery_text"`
	DeliveryFee  float64     `json:"delivery_fee"`
	ChangeFor    *float64    `json:"change_for"`
	Comment      *string     `json:"comment"`
	Rating       *int        `json:"rating"`
	OrderComment *string     `json:"order_comment"`
}

// OrderItem единица товара в заказе (универсальная модель)
type OrderItem struct {
	OrderID   int      `json:"order_id" db:"order_id"`
	ProductID int      `json:"product_id" db:"product_id"`
	Name      *string  `json:"name" db:"product_name"`
	Quantity  int      `json:"quantity" db:"quantity"`
	Price     *float64 `json:"price" db:"price"`
}

// Order полный заказ, возвращаемый клиенту
type Order struct {
	ID           int         `json:"id" db:"id"`
	OwnerID      string      `json:"owner_id" db:"owner_id"`
	Total        float64     `json:"total" db:"total"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	DateOrders   int64       `json:"date_orders"`
	Status       string      `json:"status" db:"status"`
	Name         *string     `json:"name" db:"name"`
	Address      *string     `json:"address" db:"address"`
	DeliveryType string      `json:"delivery_type" db:"delivery_type"`
	PaymentType  string      `json:"payment_type" db:"payment_type"`
	ChangeFor    *float64    `json:"change_for" db:"change_for"`
	DeliveryFee  *float64    `json:"delivery_fee" db:"delivery_fee"`
	DeliveryText *string     `json:"delivery_text" db:"delivery_text"`
	Comment      *string     `json:"comment" db:"comment"`
	Rating       *int        `json:"rating" db:"rating"`
	Items        []OrderItem `json:"items"`
	OrderComment *string     `json:"order_comment" db:"order_comment"`
}

// OrderStatusRequest используется при PATCH-запросе на обновление статуса
type OrderStatusRequest struct {
	Status string `json:"status" example:"в пути"`
}

var AllowedOrderStatuses = map[string]bool{
	"новый":      true,
	"принят":     true,
	"собирается": true,
	"отклонен":   true,
	"готов":      true,
	"в пути":     true,
	"доставлен":  true,
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		*Alias
		DateOrders int64 `json:"date_orders"`
	}{
		Alias:      (*Alias)(o),
		DateOrders: o.CreatedAt.UnixMilli(),
	})
}
