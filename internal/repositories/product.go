package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
	Delete(id int) error
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
