package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByID(id int) (*models.User, error)
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	query := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Username, user.Password, user.Role)
	return err
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE username=$1`
	err := r.db.Get(&user, query, username)
	return &user, err
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id=$1`
	err := r.db.Get(&user, query, id)
	return &user, err
}
