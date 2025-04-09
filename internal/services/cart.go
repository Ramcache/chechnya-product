package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"fmt"
)

type CartService struct {
	repo        repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(repo repositories.CartRepository, productRepo repositories.ProductRepository) *CartService {
	return &CartService{repo: repo, productRepo: productRepo}
}

func (s *CartService) AddToCart(userID, productID, quantity int) error {
	// Получаем продукт
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	// Смотрим, сколько уже в корзине
	existingItem, err := s.repo.GetCartItem(userID, productID)
	if err != nil {
		return err
	}

	totalQty := quantity
	if existingItem != nil {
		totalQty += existingItem.Quantity
	}

	// Проверяем остаток
	if totalQty > product.Stock {
		return fmt.Errorf("на складе доступно только %d шт.", product.Stock-existingItem.Quantity)
	}

	// Добавляем
	return s.repo.AddItem(userID, productID, quantity)
}

func (s *CartService) GetCart(userID int) ([]models.CartItem, error) {
	return s.repo.GetCartItems(userID)
}

func (s *CartService) UpdateItem(userID, productID, quantity int) error {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}
	if quantity > product.Stock {
		return fmt.Errorf("на складе только %d шт.", product.Stock)
	}
	return s.repo.UpdateQuantity(userID, productID, quantity)
}

func (s *CartService) DeleteItem(userID, productID int) error {
	return s.repo.DeleteItem(userID, productID)
}
