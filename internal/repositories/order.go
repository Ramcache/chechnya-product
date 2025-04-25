package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(ownerID string, total float64) (int, error)
	GetByOwnerID(ownerID string) ([]models.Order, error)
	GetAll() ([]models.Order, error)
	UpdateStatus(orderID int, status string) error
	GetByID(orderID int) (*models.Order, error)
	GetOrderItems(orderID int) ([]models.OrderItem, error)
	GetWithItemsByOwnerID(ownerID string) ([]models.OrderWithItems, error)
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

func (r *OrderRepo) UpdateStatus(orderID int, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, orderID)
	return err
}

func (r *OrderRepo) GetByID(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.db.Get(&order, `SELECT id, owner_id, total, created_at, status FROM orders WHERE id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) GetOrderItems(orderID int) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.Select(&items, `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = $1
	`, orderID)
	return items, err
}

func (r *OrderRepo) GetWithItemsByOwnerID(ownerID string) ([]models.OrderWithItems, error) {
	var orders []models.OrderWithItems

	// 1. Получаем заказы пользователя
	err := r.db.Select(&orders, `
		SELECT id, owner_id, total, status, created_at
		FROM orders
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`, ownerID)
	if err != nil {
		return nil, err
	}

	// 2. Для каждого заказа — подтягиваем товары
	for i := range orders {
		var items []models.OrderItemFull
		err := r.db.Select(&items, `
			SELECT p.id AS product_id, p.name AS product_name, oi.quantity, p.price
			FROM order_items oi
			JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = $1
		`, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}
