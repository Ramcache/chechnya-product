package repositories

import (
	"chechnya-product/internal/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Интерфейс для работы с пользователями в БД
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByOwnerID(ownerID string) (*models.User, error)
	FindByPhoneOrEmail(identifier string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	UpdateAddress(userID int, address *string) error
	GetAddress(userID int) (*string, error)
	ClearAddress(userID int) error
	GetUsernameByID(id string) (string, error)
}

// Репозиторий пользователей
type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Создать нового пользователя
func (r *UserRepo) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, phone, password_hash, role, is_verified, owner_id)
	          VALUES (:username, :email, :phone, :password_hash, :role, :is_verified, :owner_id)
	          RETURNING id, created_at`
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return fmt.Errorf("Не удалось создать пользователя: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.CreatedAt); err != nil {
			return fmt.Errorf("Не удалось получить данные пользователя: %w", err)
		}
	}
	return nil
}

// Получить пользователя по ID
func (r *UserRepo) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

// Получить пользователя по телефону
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

// Получить пользователя по email
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

// Получить пользователя по username
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

// Получить пользователя по owner_id
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

// Поиск пользователя по телефону, email или username
func (r *UserRepo) FindByPhoneOrEmail(identifier string) (*models.User, error) {
	if strings.Contains(identifier, "@") {
		return r.GetByEmail(identifier)
	}
	if strings.HasPrefix(identifier, "+") {
		return r.GetByPhone(identifier)
	}
	return r.GetByUsername(identifier)
}

func (r *UserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, `
		SELECT id, username, email, phone, role, is_verified, owner_id, created_at, password_hash
		FROM users
		ORDER BY created_at DESC
	`)
	return users, err
}

func (r *UserRepo) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, `
		SELECT id, username, email, phone, role, is_verified, owner_id, created_at, password_hash
		FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateAddress(userID int, address *string) error {
	query := `UPDATE users SET address = $1 WHERE id = $2`
	_, err := r.db.Exec(query, address, userID)
	if err != nil {
		return fmt.Errorf("не удалось обновить адрес: %w", err)
	}
	return nil
}

func (r *UserRepo) GetAddress(userID int) (*string, error) {
	var address *string
	err := r.db.Get(&address, "SELECT address FROM users WHERE id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить адрес: %w", err)
	}
	return address, nil
}

func (r *UserRepo) ClearAddress(userID int) error {
	query := `UPDATE users SET address = NULL WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("не удалось очистить адрес: %w", err)
	}
	return nil
}

func (r *UserRepo) GetUsernameByID(id string) (string, error) {
	var username string
	err := r.db.Get(&username, `SELECT username FROM users WHERE id = $1`, id)
	return username, err
}
