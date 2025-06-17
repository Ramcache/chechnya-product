package models

import "time"

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        *string   `db:"email" json:"email"`
	Phone        string    `db:"phone" json:"phone"`
	Role         UserRole  `db:"role" json:"role"`
	IsVerified   bool      `db:"is_verified" json:"is_verified"`
	OwnerID      string    `db:"owner_id" json:"owner_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Address      *string   `db:"address" json:"address"`
}
