package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type ProductServiceInterface interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.ProductResponse, error)
	AddProduct(product *models.Product) error
	UpdateProduct(id int, product *models.Product) (*models.ProductResponse, error)
	DeleteProduct(id int) error
	GetFiltered(
		search, category string,
		minPrice, maxPrice float64,
		limit, offset int,
		sort string,
		availability *bool,
	) ([]models.ProductResponse, error)
	AddProductsBulk(products []models.Product) ([]models.ProductResponse, error)
}

type ProductService struct {
	repo   repositories.ProductRepository
	logger *zap.Logger
}

func NewProductService(repo repositories.ProductRepository, logger *zap.Logger) *ProductService {
	return &ProductService{repo: repo, logger: logger}
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

func (s *ProductService) GetByID(id int) (*models.ProductResponse, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	categoryName, err := s.repo.GetCategoryNameByID(product.CategoryID)
	if err != nil {
		categoryName = ""
	}

	return &models.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Availability: product.Availability,
		CategoryID:   product.CategoryID,
		CategoryName: categoryName,
		Url:          product.Url,
	}, nil
}

func (s *ProductService) AddProduct(product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id int, product *models.Product) (*models.ProductResponse, error) {
	if err := validateProduct(product); err != nil {
		return nil, err
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := s.repo.UpdateTx(tx, id, product); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	updated, err := s.repo.GetByIDTx(tx, id)
	if err != nil || updated == nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to fetch updated product")
	}

	categoryName, err := s.repo.GetCategoryNameByIDTx(tx, updated.CategoryID)
	if err != nil {
		categoryName = ""
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		Availability: product.Availability,
		CategoryID:   product.CategoryID,
		CategoryName: categoryName,
		Url:          product.Url,
	}, nil

}

func validateProduct(p *models.Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("product name is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("product price must be positive")
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
	availability *bool,
) ([]models.ProductResponse, error) {
	products, err := s.repo.GetFiltered(search, category, minPrice, maxPrice, limit, offset, sort, availability)
	if err != nil {
		return nil, err
	}

	var result []models.ProductResponse
	for _, p := range products {
		categoryName, err := s.repo.GetCategoryNameByID(p.CategoryID)
		if err != nil {
			categoryName = ""
		}
		result = append(result, models.ProductResponse{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Availability: p.Availability,
			CategoryID:   p.CategoryID,
			CategoryName: categoryName,
			Url:          p.Url,
		})
	}
	return result, nil
}

func (s *ProductService) AddProductsBulk(products []models.Product) ([]models.ProductResponse, error) {
	var responses []models.ProductResponse

	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	for _, p := range products {
		if err := validateProduct(&p); err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		existing, err := s.repo.GetByNameTx(tx, p.Name)
		if err == nil && existing != nil {
			existing.Availability = p.Availability
			if err := s.repo.UpdateAvailabilityTx(tx, existing.ID, p.Availability); err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			updated, err := s.repo.GetByIDTx(tx, existing.ID)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			categoryName, _ := s.repo.GetCategoryNameByIDTx(tx, updated.CategoryID)
			responses = append(responses, models.ProductResponse{
				ID:           updated.ID,
				Name:         updated.Name,
				Description:  updated.Description,
				Price:        updated.Price,
				Availability: updated.Availability,
				CategoryID:   updated.CategoryID,
				CategoryName: categoryName,
				Url:          updated.Url,
			})
			continue
		}

		if err := s.repo.CreateTx(tx, &p); err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		categoryName, _ := s.repo.GetCategoryNameByIDTx(tx, p.CategoryID)
		responses = append(responses, models.ProductResponse{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Availability: p.Availability,
			CategoryID:   p.CategoryID,
			CategoryName: categoryName,
			Url:          p.Url,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, nil
}

func (s *ProductService) GetAverageRating(productID int) (float64, error) {
	rating, err := s.repo.GetAverageRating(productID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch average rating: %w", err)
	}
	return rating, nil
}
