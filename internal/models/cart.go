package models

import "time"

type Cart struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type CartItem struct {
	ID        int `db:"id"`
	CartID    int `db:"cart_id"`
	ProductID int `db:"product_id"`
	Quantity  int `db:"quantity"`
}
