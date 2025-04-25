package repositories

import (
	"chechnya-product/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	Create(name string, sortOrder int) error
	Update(id int, name string, sortOrder int) error
	Delete(id int) error
	BeginTx() (*sqlx.Tx, error)
	GetByNameTx(tx *sqlx.Tx, name string) (*models.Category, error)
	CreateReturningTx(tx *sqlx.Tx, name string, sortOrder int) (*models.Category, error)
	PartialUpdate(id int, name *string, sortOrder *int) error
}

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Select(&categories, `SELECT * FROM categories ORDER BY sort_order`)
	return categories, err
}

func (r *CategoryRepo) Create(name string, sortOrder int) error {
	_, err := r.db.Exec(`INSERT INTO categories (name, sort_order) VALUES ($1, $2)`, name, sortOrder)
	return err
}

func (r *CategoryRepo) Update(id int, name string, sortOrder int) error {
	_, err := r.db.Exec(`UPDATE categories SET name = $1, sort_order = $2 WHERE id = $3`, name, sortOrder, id)
	return err
}

func (r *CategoryRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}

func (r *CategoryRepo) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *CategoryRepo) GetByNameTx(tx *sqlx.Tx, name string) (*models.Category, error) {
	var cat models.Category
	err := tx.Get(&cat, `SELECT * FROM categories WHERE name = $1`, name)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepo) CreateReturningTx(tx *sqlx.Tx, name string, sortOrder int) (*models.Category, error) {
	var cat models.Category
	err := tx.Get(&cat, `
		INSERT INTO categories (name, sort_order)
		VALUES ($1, $2)
		RETURNING id, name, sort_order
	`, name, sortOrder)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepo) PartialUpdate(id int, name *string, sortOrder *int) error {
	if name == nil && sortOrder == nil {
		return fmt.Errorf("nothing to update: at least one field must be provided")
	}

	query := "UPDATE categories SET "
	args := []interface{}{}
	setParts := []string{}
	i := 1

	if name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", i))
		args = append(args, *name)
		i++
	}

	if sortOrder != nil {
		setParts = append(setParts, fmt.Sprintf("sort_order = $%d", i))
		args = append(args, *sortOrder)
		i++
	}

	query += strings.Join(setParts, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}
