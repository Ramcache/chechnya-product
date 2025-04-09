package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"errors"
)

type OrderService struct {
	cartRepo  repositories.CartRepository
	orderRepo repositories.OrderRepository
}

func NewOrderService(cartRepo repositories.CartRepository, orderRepo repositories.OrderRepository) *OrderService {
	return &OrderService{cartRepo: cartRepo, orderRepo: orderRepo}
}

func (s *OrderService) PlaceOrder(userID int) error {
	// Получаем товары из корзины
	items, err := s.cartRepo.GetCartItems(userID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.New("корзина пуста")
	}

	// Считаем общую сумму (в будущем можно загрузить цену из продуктов)
	var total float64
	for _, item := range items {
		// В реальном магазине — брать цену из БД
		total += float64(item.Quantity) * 100.0 // допустим 100₽ за штуку временно
	}

	// Создаём заказ
	err = s.orderRepo.CreateOrder(userID, total)
	if err != nil {
		return err
	}

	// Очищаем корзину
	return s.cartRepo.ClearCart(userID)
}

func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}
