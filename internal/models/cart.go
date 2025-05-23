package models

type CartItem struct {
	ID        int `json:"id" db:"id"`
	CartID    int `json:"cart_id" db:"cart_id"`
	ProductID int `json:"product_id" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
}
type CartItemResponse struct {
	ProductID int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

type CartBulkResponse struct {
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}
