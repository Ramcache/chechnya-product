package repositories

import (
	"chechnya-product/internal/models"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	AddItem(ownerID string, productID, quantity int) error
	GetCartItems(ownerID string) ([]models.CartItem, error)
	GetCartItem(ownerID string, productID int) (*models.CartItem, error)
	UpdateQuantity(ownerID string, productID, quantity int) error
	DeleteItem(ownerID string, productID int) error
	ClearCart(ownerID string) error
	TransferOwnership(from, to string) error
}

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) AddItem(ownerID string, productID, quantity int) error {
	_, err := r.db.Exec(`
		INSERT INTO cart_items (owner_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (owner_id, product_id) DO UPDATE
		SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, ownerID, productID, quantity)
	return err
}

func (r *CartRepo) GetCartItems(ownerID string) ([]models.CartItem, error) {
	const query = `
		SELECT id, product_id, quantity
		FROM cart_items
		WHERE owner_id = $1
	`
	var items []models.CartItem
	err := r.db.Select(&items, query, ownerID)
	return items, err
}

func (r *CartRepo) GetCartItem(ownerID string, productID int) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Get(&item, `
		SELECT id, product_id, quantity
		FROM cart_items
		WHERE owner_id = $1 AND product_id = $2
	`, ownerID, productID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}

func (r *CartRepo) UpdateQuantity(ownerID string, productID, quantity int) error {
	_, err := r.db.Exec(`
		UPDATE cart_items
		SET quantity = $1
		WHERE owner_id = $2 AND product_id = $3
	`, quantity, ownerID, productID)
	return err
}

func (r *CartRepo) DeleteItem(ownerID string, productID int) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_items
		WHERE owner_id = $1 AND product_id = $2
	`, ownerID, productID)
	return err
}

func (r *CartRepo) ClearCart(ownerID string) error {
	_, err := r.db.Exec(`
	DELETE FROM cart_items
	WHERE owner_id = $1
	`, ownerID)
	return err
}

func (r *CartRepo) TransferOwnership(from, to string) error {
	_, err := r.db.Exec(`UPDATE cart_items SET owner_id = $2 WHERE owner_id = $1`, from, to)
	return err
}
