package repositories

import (
	"chechnya-product/internal/models"
	"fmt"
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
	CreateFullOrder(ownerID string, req models.PlaceOrderRequest) (int, error)
}

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

const orderFields = `
	id, owner_id, total, created_at, status,
	name, address, delivery_type, payment_type, change_for,
	delivery_fee, delivery_text, frontend_created_at
`

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
	query := fmt.Sprintf("SELECT %s FROM orders WHERE owner_id = $1 ORDER BY created_at DESC", orderFields)
	err := r.db.Select(&orders, query, ownerID)

	return orders, err
}

func (r *OrderRepo) GetAll() ([]models.Order, error) {
	var orders []models.Order
	query := fmt.Sprintf("SELECT %s FROM orders ORDER BY created_at DESC", orderFields)
	err := r.db.Select(&orders, query)

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

func (r *OrderRepo) CreateFullOrder(ownerID string, req models.PlaceOrderRequest) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (
			owner_id, total, status, created_at, delivery_type, payment_type, 
			change_for, name, address, delivery_fee, delivery_text, frontend_created_at
		) VALUES (
			$1, $2, $3, NOW(), $4, $5, $6, $7, $8, $9, $10, $11
		)
		RETURNING id
	`, ownerID, req.Total, req.Status, req.DeliveryType, req.PaymentType, req.ChangeFor,
		req.Name, req.Address, req.DeliveryFee, req.DeliveryText, req.CreatedAt).Scan(&orderID)

	if err != nil {
		return 0, err
	}

	for _, item := range req.Items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, product_name, price)
			VALUES ($1, $2, $3, $4, $5)
		`, orderID, item.ID, item.Quantity, item.Name, item.Price)

		if err != nil {
			return 0, err
		}
	}

	return orderID, tx.Commit()
}
