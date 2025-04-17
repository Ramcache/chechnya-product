package services

import (
	"chechnya-product/internal/repositories"
	"errors"
	"fmt"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductOutOfStock   = errors.New("not enough stock available")
	ErrInvalidCartQuantity = errors.New("invalid quantity for cart operation")
)

type CartService struct {
	repo        repositories.CartRepository
	productRepo repositories.ProductRepository
}

type CartItemResponse struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

func NewCartService(repo repositories.CartRepository, productRepo repositories.ProductRepository) *CartService {
	return &CartService{repo: repo, productRepo: productRepo}
}

type CartServiceInterface interface {
	AddToCart(userID, productID, quantity int) error
	GetCart(userID int) ([]CartItemResponse, error)
	UpdateItem(userID, productID, quantity int) error
	DeleteItem(userID, productID int) error
	ClearCart(userID int) error
	Checkout(userID int) error
}

func (s *CartService) AddToCart(userID, productID, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidCartQuantity
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("fetch product: %w", err)
	}
	if product == nil {
		return ErrProductNotFound
	}

	item, err := s.repo.GetCartItem(userID, productID)
	if err != nil {
		return fmt.Errorf("get cart item: %w", err)
	}

	currentQty := 0
	if item != nil {
		currentQty = item.Quantity
	}
	if currentQty+quantity > product.Stock {
		return ErrProductOutOfStock
	}

	return s.repo.AddItem(userID, productID, quantity)
}

func (s *CartService) GetCart(userID int) ([]CartItemResponse, error) {
	items, err := s.repo.GetCartItems(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}

	var result []CartItemResponse
	for _, item := range items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil || product == nil {
			continue // или логируем, или пропускаем
		}

		response := CartItemResponse{
			ProductID: item.ProductID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  item.Quantity,
			Total:     product.Price * float64(item.Quantity),
		}
		result = append(result, response)
	}
	return result, nil
}

func (s *CartService) UpdateItem(userID, productID, quantity int) error {
	if quantity < 0 {
		return ErrInvalidCartQuantity
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("fetch product: %w", err)
	}
	if product == nil {
		return ErrProductNotFound
	}
	if quantity > product.Stock {
		return ErrProductOutOfStock
	}

	return s.repo.UpdateQuantity(userID, productID, quantity)
}

func (s *CartService) DeleteItem(userID, productID int) error {
	return s.repo.DeleteItem(userID, productID)
}

func (s *CartService) ClearCart(userID int) error {
	return s.repo.ClearCart(userID)
}

func (s *CartService) Checkout(userID int) error {
	// в будущем: списание товара, создание заказа и т.д.
	return s.repo.Checkout(userID)
}
