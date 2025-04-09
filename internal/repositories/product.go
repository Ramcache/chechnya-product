package repositories

import (
	"chechnya-product/internal/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
	Delete(id int) error
	Update(id int, product *models.Product) error
	GetByID(id int) (*models.Product, error)
	DecreaseStock(id int, quantity int) error
	GetFiltered(search, category string, limit, offset int, sort string) ([]models.Product, error)
	GetCategories() ([]string, error)
}
type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Select(&products, "SELECT * FROM products")
	return products, err
}

func (r *ProductRepo) Create(product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price, stock)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock)
	return err
}

func (r *ProductRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepo) Update(id int, product *models.Product) error {
	query := `
	UPDATE products
	SET name = $1, description = $2, price = $3, stock = $4
	WHERE id = $5
	`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock, id)
	return err
}

func (r *ProductRepo) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.Get(&p, "SELECT * FROM products WHERE id=$1", id)
	return &p, err
}

func (r *ProductRepo) DecreaseStock(id int, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2 AND stock >= $1
	`
	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("недостаточно товара на складе")
	}
	return nil
}

func (r *ProductRepo) GetFiltered(search, category string, limit, offset int, sort string) ([]models.Product, error) {
	query := `SELECT * FROM products WHERE 1=1`
	args := []interface{}{}
	i := 1

	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", i, i)
		args = append(args, "%"+strings.ToLower(search)+"%")
		i++
	}
	if category != "" {
		query += fmt.Sprintf(" AND LOWER(category) = $%d", i)
		args = append(args, strings.ToLower(category))
		i++
	}

	// 🔒 Безопасная сортировка (только whitelist)
	sortMap := map[string]string{
		"price_asc":  "price ASC",
		"price_desc": "price DESC",
		"name_asc":   "name ASC",
		"name_desc":  "name DESC",
		"stock_asc":  "stock ASC",
		"stock_desc": "stock DESC",
	}

	if orderBy, ok := sortMap[sort]; ok {
		query += " ORDER BY " + orderBy
	} else {
		query += " ORDER BY id DESC" // default
	}

	// LIMIT и OFFSET всегда последние, безопасно через args
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, limit)
		i++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", i)
		args = append(args, offset)
	}

	var products []models.Product
	err := r.db.Select(&products, query, args...)
	return products, err
}

func (r *ProductRepo) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Select(&categories, "SELECT DISTINCT category FROM products ORDER BY category ASC")
	return categories, err
}
