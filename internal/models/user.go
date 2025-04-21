package models

import "time"

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        *string   `db:"email"`
	Phone        string    `db:"phone"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	IsVerified   bool      `db:"is_verified"`
	OwnerID      string    `db:"owner_id"`
	CreatedAt    time.Time `db:"created_at"`
}
