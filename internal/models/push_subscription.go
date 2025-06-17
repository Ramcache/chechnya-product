package models

import "time"

type PushSubscription struct {
	ID        int       `db:"id"`
	Endpoint  string    `db:"endpoint"`
	P256DH    string    `db:"p256dh"`
	Auth      string    `db:"auth"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
