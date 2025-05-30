package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/ws"
	"errors"
	"fmt"
)

type OrderServiceInterface interface {
	PlaceOrder(ownerID string, req models.PlaceOrderRequest) (*models.Order, error)
	GetOrders(ownerID string) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateStatus(orderID int, status string) error
	RepeatOrder(orderID int, ownerID string) error
	GetOrderHistory(ownerID string) ([]models.Order, error)
	DeleteOrder(orderID int) error
	GetOrderByID(orderID int) (*models.Order, error)
	UpdateReview(orderID int, comment *string, rating *int) error
	AddReview(orderID int, comment *string, rating *int) error
	GetByOrderReviewID(orderID int) (*models.OrderReview, error)
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

func (s *OrderService) PlaceOrder(ownerID string, req models.PlaceOrderRequest) (*models.Order, error) {
	// 1. Считаем сумму заказа
	var total float64
	for _, item := range req.Items {
		if item.Price != nil {
			total += *item.Price * float64(item.Quantity)
		}
	}
	total += req.DeliveryFee

	// 2. Создаём заказ
	orderID, err := s.orderRepo.CreateFullOrder(ownerID, req, total)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 3. Очищаем корзину
	if err := s.cartRepo.ClearCart(ownerID); err != nil {
		return nil, fmt.Errorf("order created but failed to clear cart: %w", err)
	}

	// 4. Получаем заказ
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order created, but failed to fetch: %w", err)
	}

	// 5. Получаем товары заказа
	items, err := s.orderRepo.GetOrderItems(orderID)
	if err != nil {
		return nil, fmt.Errorf("order created, but failed to fetch items: %w", err)
	}
	order.Items = items

	// 6. WebSocket
	if s.hub != nil {
		s.hub.BroadcastNewOrder(*order)
	}

	return order, nil
}

func (s *OrderService) GetOrders(ownerID string) ([]models.Order, error) {
	orders, err := s.orderRepo.GetWithItemsByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []models.Order{}, nil
	}
	return orders, nil
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllWithItems()
}

func (s *OrderService) UpdateStatus(orderID int, status string) error {
	if !models.AllowedOrderStatuses[status] {
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

func (s *OrderService) GetOrderHistory(ownerID string) ([]models.Order, error) {
	return s.orderRepo.GetWithItemsByOwnerID(ownerID)
}

func (s *OrderService) DeleteOrder(orderID int) error {
	return s.orderRepo.DeleteOrder(orderID)
}

func (s *OrderService) GetOrderByID(orderID int) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.orderRepo.GetOrderItems(orderID)
	if err != nil {
		return nil, err
	}

	order.Items = items
	return order, nil
}

func (s *OrderService) UpdateReview(orderID int, comment *string, rating *int) error {
	if rating != nil && (*rating < 1 || *rating > 5) {
		return fmt.Errorf("рейтинг должен быть от 1 до 5")
	}
	return s.orderRepo.UpdateReview(orderID, comment, rating)
}

func (s *OrderService) Add(orderID int, comment *string, rating *int) error {
	if rating != nil && (*rating < 1 || *rating > 5) {
		return fmt.Errorf("рейтинг должен быть от 1 до 5")
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("заказ не найден")
	}
	if order.Status != "доставлен" {
		return fmt.Errorf("оставлять отзыв можно только после доставки")
	}

	return s.orderRepo.AddReview(orderID, comment, rating)
}

func (s *OrderService) GetByOrderReviewID(orderID int) (*models.OrderReview, error) {
	return s.orderRepo.GetReviewByOrderID(orderID)
}
