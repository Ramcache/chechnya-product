package models

type Category struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	SortOrder int    `json:"sort_order" db:"sort_order"` // 👈 это важно для sqlx
}
