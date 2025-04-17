package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"errors"
	"fmt"
)

type OrderService struct {
	cartRepo    repositories.CartRepository
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
}

func NewOrderService(cartRepo repositories.CartRepository, orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository) *OrderService {
	return &OrderService{
		cartRepo:    cartRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) PlaceOrder(userID int) error {
	items, err := s.cartRepo.GetCartItems(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(items) == 0 {
		return errors.New("cart is empty")
	}

	var total float64
	for _, item := range items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}
		if product.Stock < item.Quantity {
			return fmt.Errorf("not enough stock for \"%s\"", product.Name)
		}

		total += float64(item.Quantity) * product.Price

		if err := s.productRepo.DecreaseStock(item.ProductID, item.Quantity); err != nil {
			return fmt.Errorf("failed to decrease stock for %d: %w", item.ProductID, err)
		}
	}

	if err := s.orderRepo.CreateOrder(userID, total); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := s.cartRepo.ClearCart(userID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}
