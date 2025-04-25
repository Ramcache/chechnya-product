package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	Create(name string, sortOrder int) error
	Update(id int, name string, sortOrder int) error
	Delete(id int) error
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
