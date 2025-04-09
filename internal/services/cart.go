package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
)

type CartService struct {
	repo repositories.CartRepository
}

func NewCartService(repo repositories.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddToCart(userID, productID, quantity int) error {
	return s.repo.AddItem(userID, productID, quantity)
}

func (s *CartService) GetCart(userID int) ([]models.CartItem, error) {
	return s.repo.GetCartItems(userID)
}
