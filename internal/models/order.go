package models

// PlaceOrderRequest запрос от клиента для оформления заказа
type PlaceOrderRequest struct {
	Name         *string     `json:"name"`
	Address      *string     `json:"address"`
	Items        []OrderItem `json:"items"`
	Total        float64     `json:"total"`
	PaymentType  string      `json:"paymentType"`
	Status       string      `json:"status"`
	DeliveryType string      `json:"deliveryType"`
	CreatedAt    int64       `json:"createdAt"`
	DeliveryText string      `json:"deliveryText"`
	DeliveryFee  float64     `json:"deliveryFee"`
	ChangeFor    *float64    `json:"changeFor"`
}

// OrderItem единица товара в заказе (универсальная модель)
type OrderItem struct {
	OrderID   int      `json:"order_id" db:"order_id"`
	ProductID int      `json:"product_id" db:"product_id"`
	Name      *string  `json:"name,omitempty" db:"product_name"`
	Quantity  int      `json:"quantity" db:"quantity"`
	Price     *float64 `json:"price" db:"price"`
}

// Order полный заказ, возвращаемый клиенту
type Order struct {
	ID                int         `json:"id" db:"id"`
	OwnerID           string      `json:"owner_id" db:"owner_id"`
	Total             float64     `json:"total" db:"total"`
	CreatedAt         string      `json:"created_at" db:"created_at"`
	Status            string      `json:"status" db:"status"`
	Name              *string     `json:"name" db:"name"`
	Address           *string     `json:"address" db:"address"`
	DeliveryType      string      `json:"delivery_type" db:"delivery_type"`
	PaymentType       string      `json:"payment_type" db:"payment_type"`
	ChangeFor         *float64    `json:"change_for" db:"change_for"`
	DeliveryFee       *float64    `json:"delivery_fee" db:"delivery_fee"`
	DeliveryText      *string     `json:"delivery_text" db:"delivery_text"`
	FrontendCreatedAt *int64      `json:"frontend_created_at" db:"frontend_created_at"`
	Items             []OrderItem `json:"items"` // всегда с товарами
}

// OrderStatusRequest используется при PATCH-запросе на обновление статуса
type OrderStatusRequest struct {
	Status string `json:"status" example:"в пути"`
}
