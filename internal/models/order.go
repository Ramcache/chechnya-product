package models

import "time"

type Order struct {
	ID        int       `db:"id"`
	OwnerID   string    `db:"owner_id"`
	Total     float64   `db:"total"`
	CreatedAt time.Time `db:"created_at"`
}
