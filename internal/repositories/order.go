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
	DeleteOrder(orderID int) error
	AddReview(orderID int, comment *string, rating *int, userID int) error
	GetReviewByOrderID(orderID int) (*models.OrderReview, error)
	GetAllOrderReviews() ([]models.OrderReview, error)
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
	delivery_fee, delivery_text, order_comment
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
	err := r.db.Get(&order, `
    SELECT id, owner_id, total, created_at, status, name, address,
       delivery_type, payment_type, change_for, delivery_fee, delivery_text,
       order_comment
	FROM orders 
	WHERE id = $1

`, orderID)

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
		owner_id, total, status, created_at,
		delivery_type, payment_type, change_for,
		name, address, delivery_fee, delivery_text,
		order_comment
	) VALUES (
		$1, $2, $3, NOW(),
		$4, $5, $6,
		$7, $8, $9, $10,
		$11
	)
	RETURNING id
`,
		ownerID,
		total,
		req.Status,
		req.DeliveryType,
		req.PaymentType,
		req.ChangeFor,
		req.Name,
		req.Address,
		req.DeliveryFee,
		req.DeliveryText,
		req.OrderComment,
	).Scan(&orderID)

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

func (r *OrderRepo) DeleteOrder(orderID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Удалим сначала товары заказа (связанные строки)
	if _, err := tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, orderID); err != nil {
		tx.Rollback()
		return err
	}

	// Удалим сам заказ
	res, err := tx.Exec(`DELETE FROM orders WHERE id = $1`, orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("order with id %d not found", orderID)
	}

	return tx.Commit()
}

func (r *OrderRepo) AddReview(orderID int, comment *string, rating *int, userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO order_reviews (order_id, comment, rating, user_id)
		VALUES ($1, $2, $3, $4)
	`, orderID, comment, rating, userID)
	return err
}

func (r *OrderRepo) GetReviewByOrderID(orderID int) (*models.OrderReview, error) {
	var review models.OrderReview
	err := r.db.Get(&review, `SELECT * FROM order_reviews WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *OrderRepo) GetAllOrderReviews() ([]models.OrderReview, error) {
	var reviews []models.OrderReview
	err := r.db.Select(&reviews, `
		SELECT orr.id, orr.order_id, orr.user_id, u.username, orr.rating, orr.comment, orr.created_at
		FROM order_reviews orr
		LEFT JOIN users u ON orr.user_id = u.id
		ORDER BY orr.created_at DESC
	`)
	return reviews, err
}
