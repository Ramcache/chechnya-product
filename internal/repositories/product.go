package repositories

import (
	"chechnya-product/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type ProductRepository interface {
	BeginTx() (*sqlx.Tx, error)
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
	Delete(id int) error
	Update(id int, product *models.Product) error
	GetByID(id int) (*models.Product, error)
	DecreaseStock(id int, quantity int) error
	GetFiltered(search, category string, minPrice, maxPrice float64, limit, offset int, sort string) ([]models.Product, error)
	GetCategories() ([]string, error)
	GetCategoryNameByID(categoryID int) (string, error)
	GetByName(name string) (*models.Product, error)
	AddStock(id, additionalStock int) error
	CreateTx(tx *sqlx.Tx, p *models.Product) error
	AddStockTx(tx *sqlx.Tx, id, amount int) error
	GetByNameTx(tx *sqlx.Tx, name string) (*models.Product, error)
	GetByIDTx(tx *sqlx.Tx, id int) (*models.Product, error)
	GetCategoryNameByIDTx(tx *sqlx.Tx, categoryID int) (string, error)
	UpdateTx(tx *sqlx.Tx, id int, p *models.Product) error
}

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *ProductRepo) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Select(&products, `SELECT * FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}
	return products, nil
}

func (r *ProductRepo) Create(product *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepo) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *ProductRepo) Update(id int, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, category_id = $5
		WHERE id = $6
	`
	result, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *ProductRepo) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.Get(&p, `SELECT * FROM products WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return &p, nil
}

func (r *ProductRepo) DecreaseStock(id int, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2 AND stock >= $1
	`
	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		return fmt.Errorf("failed to decrease stock: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not enough stock or product not found")
	}
	return nil
}

func (r *ProductRepo) GetFiltered(
	search, category string,
	minPrice, maxPrice float64,
	limit, offset int,
	sort string,
) ([]models.Product, error) {
	query := `SELECT * FROM products WHERE 1=1`
	args := []interface{}{}
	i := 1

	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", i, i)
		args = append(args, "%"+strings.ToLower(search)+"%")
		i++
	}
	if category != "" {
		query += fmt.Sprintf(" AND category_id = $%d", i)
		categoryID, err := strconv.Atoi(category) // если нужно из строки
		if err != nil {
			return nil, fmt.Errorf("invalid category ID")
		}
		args = append(args, categoryID)
		i++
	}

	if minPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", i)
		args = append(args, minPrice)
		i++
	}
	if maxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", i)
		args = append(args, maxPrice)
		i++
	}

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
		query += " ORDER BY id DESC"
	}

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
	if err != nil {
		return nil, fmt.Errorf("failed to filter products: %w", err)
	}
	return products, nil
}

func (r *ProductRepo) GetCategories() ([]string, error) {
	var names []string
	err := r.db.Select(&names, `SELECT name FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return names, nil
}

func (r *ProductRepo) GetCategoryNameByID(categoryID int) (string, error) {
	var name string
	err := r.db.Get(&name, `SELECT name FROM categories WHERE id = $1`, categoryID)
	return name, err
}

func (r *ProductRepo) GetByName(name string) (*models.Product, error) {
	var product models.Product
	err := r.db.Get(&product, `SELECT * FROM products WHERE name = $1`, name)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) AddStock(id, additionalStock int) error {
	_, err := r.db.Exec(`UPDATE products SET stock = stock + $1 WHERE id = $2`, additionalStock, id)
	return err
}

func (r *ProductRepo) CreateTx(tx *sqlx.Tx, p *models.Product) error {
	query := `INSERT INTO products (name, description, price, stock, category_id)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id`
	return tx.QueryRow(query, p.Name, p.Description, p.Price, p.Stock, p.CategoryID).Scan(&p.ID)
}

func (r *ProductRepo) AddStockTx(tx *sqlx.Tx, id, amount int) error {
	_, err := tx.Exec(`UPDATE products SET stock = stock + $1 WHERE id = $2`, amount, id)
	return err
}

func (r *ProductRepo) GetByNameTx(tx *sqlx.Tx, name string) (*models.Product, error) {
	var p models.Product
	err := tx.Get(&p, `SELECT * FROM products WHERE name = $1`, name)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) GetByIDTx(tx *sqlx.Tx, id int) (*models.Product, error) {
	var p models.Product
	err := tx.Get(&p, `SELECT * FROM products WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func (r *ProductRepo) GetCategoryNameByIDTx(tx *sqlx.Tx, categoryID int) (string, error) {
	var name string
	err := tx.Get(&name, `SELECT name FROM categories WHERE id = $1`, categoryID)
	return name, err
}

func (r *ProductRepo) UpdateTx(tx *sqlx.Tx, id int, p *models.Product) error {
	_, err := tx.Exec(`UPDATE products SET name=$1, description=$2, price=$3, stock=$4, category_id=$5 WHERE id=$6`,
		p.Name, p.Description, p.Price, p.Stock, p.CategoryID, id)
	return err
}
