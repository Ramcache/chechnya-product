package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"fmt"
)

type ReviewServiceInterface interface {
	AddReview(ownerID string, productID, rating int, comment string) error
	GetReviewsByProductID(productID int) ([]models.Review, error)
}

type ReviewService struct {
	repo repositories.ReviewRepository
}

func NewReviewService(repo repositories.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) AddReview(ownerID string, productID, rating int, comment string) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	review := &models.Review{
		OwnerID:   ownerID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	}
	return s.repo.Create(review)
}

func (s *ReviewService) GetReviewsByProductID(productID int) ([]models.Review, error) {
	return s.repo.GetByProductID(productID)
}
