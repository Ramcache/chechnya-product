package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/ws"
	"errors"
	"fmt"
	"log"
	"time"
)

type OrderService struct {
	cartRepo    repositories.CartRepository
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
	hub         *ws.Hub
}

func NewOrderService(
	cartRepo repositories.CartRepository,
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	hub *ws.Hub,
) *OrderService {
	return &OrderService{
		cartRepo:    cartRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		hub:         hub,
	}
}

func (s *OrderService) PlaceOrder(ownerID string) error {
	items, err := s.cartRepo.GetCartItems(ownerID)
	log.Println("[OWNER]", ownerID)

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

	orderID, err := s.orderRepo.CreateOrder(ownerID, total)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := s.cartRepo.ClearCart(ownerID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	order := models.Order{
		ID:        orderID,
		OwnerID:   ownerID,
		Total:     total,
		CreatedAt: time.Now(), // опционально: получи из БД
	}

	if s.hub != nil {
		s.hub.BroadcastNewOrder(order)
	}

	return nil
}

func (s *OrderService) GetOrders(ownerID string) ([]models.Order, error) {
	orders, err := s.orderRepo.GetByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []models.Order{}, nil
	}
	return orders, nil
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []models.Order{}, nil
	}
	return orders, nil
}
