package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"fmt"
	"strings"
)

type ProductServiceInterface interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	AddProduct(product *models.Product) error
	UpdateProduct(id int, product *models.Product) error
	DeleteProduct(id int) error
	GetFiltered(
		search, category string,
		minPrice, maxPrice float64,
		limit, offset int,
		sort string,
	) ([]models.Product, error)
}

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if products == nil {
		return []models.Product{}, nil
	}
	return products, nil
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

func (s *ProductService) AddProduct(product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id int, product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Update(id, product)
}

func validateProduct(p *models.Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("product name is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("product price must be positive")
	}
	if p.Stock < 0 {
		return fmt.Errorf("product stock cannot be negative")
	}
	return nil
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetFiltered(
	search, category string,
	minPrice, maxPrice float64,
	limit, offset int,
	sort string,
) ([]models.Product, error) {
	products, err := s.repo.GetFiltered(search, category, minPrice, maxPrice, limit, offset, sort)
	if err != nil {
		return nil, err
	}
	if products == nil {
		return []models.Product{}, nil
	}
	return products, nil
}
