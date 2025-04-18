package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(ownerID string, total float64) (int, error)
	GetByOwnerID(ownerID string) ([]models.Order, error)
	GetAll() ([]models.Order, error)
}

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(ownerID string, total float64) (int, error) {
	var orderID int
	err := r.db.QueryRow(`
		INSERT INTO orders (owner_id, total)
		VALUES ($1, $2)
		RETURNING id
	`, ownerID, total).Scan(&orderID)
	return orderID, err
}

func (r *OrderRepo) GetByOwnerID(ownerID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Select(&orders, `
		SELECT * FROM orders 
		WHERE owner_id = $1 
		ORDER BY created_at DESC
	`, ownerID)
	return orders, err
}

func (r *OrderRepo) GetAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Select(&orders, `
		SELECT * FROM orders 
		ORDER BY created_at DESC
	`)
	return orders, err
}
