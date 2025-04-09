package models

type CartItem struct {
	ID        int `json:"id" db:"id"`
	CartID    int `json:"cart_id" db:"cart_id"`
	ProductID int `json:"product_id" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
}
