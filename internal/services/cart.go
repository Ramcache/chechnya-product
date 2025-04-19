package services

import (
	"chechnya-product/internal/repositories"
	"errors"
	"fmt"
	"log"
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
	AddToCart(ownerID string, productID, quantity int) error
	GetCart(ownerID string) ([]CartItemResponse, error)
	UpdateItem(ownerID string, productID, quantity int) error
	DeleteItem(ownerID string, productID int) error
	ClearCart(ownerID string) error
	Checkout(ownerID string) error
}

func (s *CartService) AddToCart(ownerID string, productID, quantity int) error {
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

	item, err := s.repo.GetCartItem(ownerID, productID)
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

	return s.repo.AddItem(ownerID, productID, quantity)
}

func (s *CartService) GetCart(ownerID string) ([]CartItemResponse, error) {
	items, err := s.repo.GetCartItems(ownerID)
	log.Println("[OWNER]", ownerID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}

	result := make([]CartItemResponse, 0)

	for _, item := range items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil || product == nil {
			continue
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

func (s *CartService) UpdateItem(ownerID string, productID, quantity int) error {
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

	return s.repo.UpdateQuantity(ownerID, productID, quantity)
}

func (s *CartService) DeleteItem(ownerID string, productID int) error {
	return s.repo.DeleteItem(ownerID, productID)
}

func (s *CartService) ClearCart(ownerID string) error {
	return s.repo.ClearCart(ownerID)
}
