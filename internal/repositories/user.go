package repositories

import (
	"chechnya-product/internal/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByOwnerID(ownerID string) (*models.User, error)
	FindByPhoneOrEmail(identifier string) (*models.User, error)
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, phone, password_hash, role, is_verified, owner_id)
	          VALUES (:username, :email, :phone, :password_hash, :role, :is_verified, :owner_id)
	          RETURNING id, created_at`
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan user data: %w", err)
		}
	}
	return nil
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE phone = $1", phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByOwnerID(ownerID string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE owner_id = $1", ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByPhoneOrEmail(identifier string) (*models.User, error) {
	if strings.Contains(identifier, "@") {
		return r.GetByEmail(identifier)
	}
	if strings.HasPrefix(identifier, "+") {
		return r.GetByPhone(identifier)
	}
	return r.GetByUsername(identifier)
}
