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
	return &OrderService{cartRepo: cartRepo, orderRepo: orderRepo, productRepo: productRepo}
}

func (s *OrderService) PlaceOrder(userID int) error {
	items, err := s.cartRepo.GetCartItems(userID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return errors.New("корзина пуста")
	}

	var total float64
	for _, item := range items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < item.Quantity {
			return fmt.Errorf("товара \"%s\" недостаточно на складе", product.Name)
		}

		total += float64(item.Quantity) * product.Price

		// уменьшаем остатки
		err = s.productRepo.DecreaseStock(item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	err = s.orderRepo.CreateOrder(userID, total)
	if err != nil {
		return err
	}

	return s.cartRepo.ClearCart(userID)
}

func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}
