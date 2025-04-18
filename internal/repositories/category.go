package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	GetAll() ([]models.Category, error)
	Create(name string) error
	Update(id int, name string) error
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
	err := r.db.Select(&categories, `SELECT * FROM categories ORDER BY id`)
	return categories, err
}

func (r *CategoryRepo) Create(name string) error {
	_, err := r.db.Exec(`INSERT INTO categories (name) VALUES ($1)`, name)
	return err
}

func (r *CategoryRepo) Update(id int, name string) error {
	_, err := r.db.Exec(`UPDATE categories SET name = $1 WHERE id = $2`, name, id)
	return err
}

func (r *CategoryRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}
