package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	GetByProductID(productID int) ([]models.Review, error)
	Exists(ownerID string, productID int) (bool, error)
	Update(ownerID string, productID, rating int, comment string) error
	Delete(ownerID string, productID int) error
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

func (r *ReviewRepo) Exists(ownerID string, productID int) (bool, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM reviews WHERE owner_id = $1 AND product_id = $2
	`, ownerID, productID)
	return count > 0, err
}

func (r *ReviewRepo) Update(ownerID string, productID, rating int, comment string) error {
	_, err := r.db.Exec(`
		UPDATE reviews SET rating=$1, comment=$2 WHERE owner_id=$3 AND product_id=$4
	`, rating, comment, ownerID, productID)
	return err
}

func (r *ReviewRepo) Delete(ownerID string, productID int) error {
	_, err := r.db.Exec(`
		DELETE FROM reviews WHERE owner_id=$1 AND product_id=$2
	`, ownerID, productID)
	return err
}
