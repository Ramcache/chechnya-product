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
	GetWithItemsByOwnerID(ownerID string) ([]models.Order, error)
	CreateFullOrder(ownerID string, req models.PlaceOrderRequest, total float64) (int, error)
	GetAllWithItems() ([]models.Order, error)
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
	delivery_fee, delivery_text
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
	res, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, orderID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("order with id %d not found", orderID)
	}

	return nil
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
		SELECT order_id, product_id, product_name, quantity, price
		FROM order_items
		WHERE order_id = $1
	`, orderID)
	return items, err
}

func (r *OrderRepo) GetWithItemsByOwnerID(ownerID string) ([]models.Order, error) {
	var orders []models.Order

	query := fmt.Sprintf("SELECT %s FROM orders WHERE owner_id = $1 ORDER BY created_at DESC", orderFields)
	err := r.db.Select(&orders, query, ownerID)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := r.GetOrderItems(orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (r *OrderRepo) CreateFullOrder(ownerID string, req models.PlaceOrderRequest, total float64) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}

	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (
					owner_id, total, status, created_at, delivery_type, payment_type, 
					change_for, name, address, delivery_fee, delivery_text
		) VALUES (
					$1, $2, $3, NOW(), $4, $5, $6, $7, $8, $9, $10
		)

		RETURNING id
	`, ownerID, total, req.Status, req.DeliveryType, req.PaymentType, req.ChangeFor,
		req.Name, req.Address, req.DeliveryFee, req.DeliveryText).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, item := range req.Items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, product_name, price)
			VALUES ($1, $2, $3, $4, $5)
		`, orderID, item.ProductID, item.Quantity, item.Name, item.Price)

		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *OrderRepo) GetAllWithItems() ([]models.Order, error) {
	var orders []models.Order

	query := fmt.Sprintf("SELECT %s FROM orders ORDER BY created_at DESC", orderFields)
	err := r.db.Select(&orders, query)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := r.GetOrderItems(orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}
