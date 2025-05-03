package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/ws"
	"errors"
	"fmt"
	"time"
)

type OrderServiceInterface interface {
	PlaceOrder(ownerID string, req models.PlaceOrderRequest) error
	GetOrders(ownerID string) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateStatus(orderID int, status string) error
	RepeatOrder(orderID int, ownerID string) error
	GetOrderHistory(ownerID string) ([]models.OrderWithItems, error)
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

func (s *OrderService) PlaceOrder(ownerID string, req models.PlaceOrderRequest) error {
	orderID, err := s.orderRepo.CreateFullOrder(ownerID, req)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Очистим корзину после создания заказа
	if err := s.cartRepo.ClearCart(ownerID); err != nil {
		return fmt.Errorf("order created but failed to clear cart: %w", err)
	}

	// Отправим уведомление через WebSocket, если используешь
	if s.hub != nil {
		s.hub.BroadcastNewOrder(models.Order{
			ID:        orderID,
			OwnerID:   ownerID,
			Total:     req.Total,
			Status:    req.Status,
			CreatedAt: time.Now(),
		})
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

func (s *OrderService) RepeatOrder(orderID int, ownerID string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil || order.OwnerID != ownerID {
		return errors.New("invalid order")
	}

	items, err := s.orderRepo.GetOrderItems(orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.cartRepo.AddOrUpdate(ownerID, item.ProductID, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderService) GetOrderHistory(ownerID string) ([]models.OrderWithItems, error) {
	return s.orderRepo.GetWithItemsByOwnerID(ownerID)
}
