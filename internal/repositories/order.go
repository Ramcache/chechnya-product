package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(userID int, total float64) error
	GetByUserID(userID int) ([]models.Order, error)
	GetAll() ([]models.Order, error) // для админа
}

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(userID int, total float64) error {
	_, err := r.db.Exec("INSERT INTO orders (user_id, total) VALUES ($1, $2)", userID, total)
	return err
}

func (r *OrderRepo) GetByUserID(userID int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Select(&orders, "SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC", userID)
	return orders, err
}

func (r *OrderRepo) GetAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Select(&orders, "SELECT * FROM orders ORDER BY created_at DESC")
	return orders, err
}
