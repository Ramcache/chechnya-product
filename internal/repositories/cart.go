package repositories

import (
	"chechnya-product/internal/models"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	AddItem(userID, productID, quantity int) error
	GetCartItems(userID int) ([]models.CartItem, error)
	ClearCart(userID int) error
	GetCartItem(userID, productID int) (*models.CartItem, error)
	UpdateQuantity(userID, productID, quantity int) error
	DeleteItem(userID, productID int) error
	Checkout(userID int) error
}

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) AddItem(userID, productID, quantity int) error {
	cartID, err := r.getOrCreateCartID(userID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id) DO UPDATE
		SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, cartID, productID, quantity)

	return err
}

func (r *CartRepo) GetCartItems(userID int) ([]models.CartItem, error) {
	const query = `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity
		FROM cart_items ci
		JOIN carts c ON ci.cart_id = c.id
		WHERE c.user_id = $1
	`
	var items []models.CartItem
	err := r.db.Select(&items, query, userID)
	return items, err
}

func (r *CartRepo) ClearCart(userID int) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_items
		WHERE cart_id IN (SELECT id FROM carts WHERE user_id = $1)
	`, userID)
	return err
}

func (r *CartRepo) GetCartItem(userID, productID int) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Get(&item, `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity
		FROM cart_items ci
		JOIN carts c ON ci.cart_id = c.id
		WHERE c.user_id = $1 AND ci.product_id = $2
	`, userID, productID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}

func (r *CartRepo) UpdateQuantity(userID, productID, quantity int) error {
	_, err := r.db.Exec(`
		UPDATE cart_items
		SET quantity = $1
		WHERE product_id = $2 AND cart_id = (
			SELECT id FROM carts WHERE user_id = $3
		)
	`, quantity, productID, userID)
	return err
}

func (r *CartRepo) DeleteItem(userID, productID int) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_items
		WHERE product_id = $1 AND cart_id = (
			SELECT id FROM carts WHERE user_id = $2
		)
	`, productID, userID)
	return err
}

func (r *CartRepo) Checkout(userID int) error {
	// В простом варианте просто очищаем корзину
	return r.ClearCart(userID)
}

// -- внутренние хелперы --

func (r *CartRepo) getOrCreateCartID(userID int) (int, error) {
	var cartID int
	err := r.db.Get(&cartID, `
		INSERT INTO carts (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING id
	`, userID)
	return cartID, err
}

func (r *CartRepo) getCartItemID(cartID, productID int) (int, error) {
	var itemID int
	err := r.db.Get(&itemID, `
		SELECT id FROM cart_items WHERE cart_id = $1 AND product_id = $2
	`, cartID, productID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, sql.ErrNoRows
	}
	return itemID, err
}

func (r *CartRepo) insertNewItem(cartID, productID, quantity int) error {
	_, err := r.db.Exec(`
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
	`, cartID, productID, quantity)
	return err
}

func (r *CartRepo) incrementItemQuantity(itemID, quantity int) error {
	_, err := r.db.Exec(`
		UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2
	`, quantity, itemID)
	return err
}
