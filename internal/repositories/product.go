package repositories

import (
	"chechnya-product/internal/models"
	"errors"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
	Delete(id int) error
	Update(id int, product *models.Product) error
	GetByID(id int) (*models.Product, error)
	DecreaseStock(id int, quantity int) error
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
