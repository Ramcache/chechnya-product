package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	GetByProductID(productID int) ([]models.Review, error)
}
type ReviewRepo struct {
	db *sqlx.DB
}

func NewReviewRepo(db *sqlx.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(review *models.Review) error {
	_, err := r.db.Exec(`
		INSERT INTO reviews (owner_id, product_id, rating, comment)
		VALUES ($1, $2, $3, $4)
	`, review.OwnerID, review.ProductID, review.Rating, review.Comment)
	return err
}

func (r *ReviewRepo) GetByProductID(productID int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Select(&reviews, `SELECT * FROM reviews WHERE product_id = $1 ORDER BY created_at DESC`, productID)
	return reviews, err
}
