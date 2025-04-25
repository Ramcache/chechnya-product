package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"fmt"
)

type ReviewServiceInterface interface {
	AddReview(ownerID string, productID, rating int, comment string) error
	GetReviewsByProductID(productID int) ([]models.Review, error)
	UpdateReview(ownerID string, productID, rating int, comment string) error
	DeleteReview(ownerID string, productID int) error
}

type ReviewService struct {
	repo repositories.ReviewRepository
}

func NewReviewService(repo repositories.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) AddReview(ownerID string, productID, rating int, comment string) error {
	exists, err := s.repo.Exists(ownerID, productID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("you have already left a review for this product")
	}
	return s.repo.Create(&models.Review{
		OwnerID:   ownerID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	})
}

func (s *ReviewService) GetReviewsByProductID(productID int) ([]models.Review, error) {
	return s.repo.GetByProductID(productID)
}
func (s *ReviewService) UpdateReview(ownerID string, productID, rating int, comment string) error {
	return s.repo.Update(ownerID, productID, rating, comment)
}

func (s *ReviewService) DeleteReview(ownerID string, productID int) error {
	return s.repo.Delete(ownerID, productID)
}
