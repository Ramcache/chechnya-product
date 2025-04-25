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

type OrderServiceInterface interface {
	PlaceOrder(ownerID string) error
	GetOrders(ownerID string) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateStatus(orderID int, status string) error
}

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
		if product == nil {
			return fmt.Errorf("product %d not found", item.ProductID)
		}
		if !product.Availability {
			return fmt.Errorf("product \"%s\" is not available", product.Name)
		}

		total += float64(item.Quantity) * product.Price
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
		CreatedAt: time.Now(),
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

func (s *OrderService) UpdateStatus(orderID int, status string) error {
	allowed := map[string]bool{
		"в обработке": true,
		"принят":      true,
		"отклонен":    true,
		"готов":       true,
		"в пути":      true,
		"доставлен":   true,
	}
	if !allowed[status] {
		return fmt.Errorf("недопустимый статус")
	}

	err := s.orderRepo.UpdateStatus(orderID, status)
	if err != nil {
		return err
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err == nil && s.hub != nil {
		s.hub.BroadcastStatusUpdate(*order)
	}

	return nil
}
