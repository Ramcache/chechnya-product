package repositories

import (
	"chechnya-product/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	AddItem(userID int, productID int, quantity int) error
	GetCartItems(userID int) ([]models.CartItem, error)
	ClearCart(userID int) error
	GetCartItem(userID, productID int) (*models.CartItem, error)
}

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

// Добавление товара в корзину
func (r *CartRepo) AddItem(userID, productID, quantity int) error {
	var cartID int

	// 1. Найти корзину или создать
	err := r.db.Get(&cartID, "SELECT id FROM carts WHERE user_id=$1", userID)
	if err == sql.ErrNoRows {
		err = r.db.QueryRow("INSERT INTO carts (user_id) VALUES ($1) RETURNING id", userID).Scan(&cartID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 2. Проверить, есть ли уже такой товар в корзине
	var itemID int
	err = r.db.Get(&itemID, "SELECT id FROM cart_items WHERE cart_id=$1 AND product_id=$2", cartID, productID)
	if err == sql.ErrNoRows {
		// Вставить новый товар
		_, err = r.db.Exec("INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3)",
			cartID, productID, quantity)
		return err
	} else if err != nil {
		return err
	}

	// Обновить количество, если товар уже есть
	_, err = r.db.Exec("UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2", quantity, itemID)
	return err
}

// Получить все товары в корзине пользователя
func (r *CartRepo) GetCartItems(userID int) ([]models.CartItem, error) {
	var items []models.CartItem

	query := `
	SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity
	FROM cart_items ci
	JOIN carts c ON ci.cart_id = c.id
	WHERE c.user_id = $1
	`

	err := r.db.Select(&items, query, userID)
	return items, err
}

func (r *CartRepo) ClearCart(userID int) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_items WHERE cart_id IN (
			SELECT id FROM carts WHERE user_id = $1
		)`, userID)
	return err
}

func (r *CartRepo) GetCartItem(userID, productID int) (*models.CartItem, error) {
	var item models.CartItem

	query := `
	SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity
	FROM cart_items ci
	JOIN carts c ON ci.cart_id = c.id
	WHERE c.user_id = $1 AND ci.product_id = $2
	`

	err := r.db.Get(&item, query, userID, productID)
	if err == sql.ErrNoRows {
		return nil, nil // ещё нет в корзине
	}
	return &item, err
}
