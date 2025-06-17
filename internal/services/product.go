package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/utils"
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
	GetFilteredRaw(
		search, category string,
		minPrice, maxPrice float64,
		limit, offset int,
		sort string,
		availability *bool,
	) ([]models.Product, error)
	GetCategoryNameByID(id int) (string, error)
	AddProductsBulk(products []models.Product) ([]models.ProductResponse, error)
	PatchProduct(id int, updates map[string]interface{}) error
	GetPaginated(
		search, category string,
		minPrice, maxPrice float64,
		limit, offset int,
		sort string,
		availability *bool,
	) ([]models.ProductResponse, int, error)
	CountFiltered(
		search, category string,
		minPrice, maxPrice float64,
		availability *bool,
	) (int, error)
}

type ProductService struct {
	repo   repositories.ProductRepository
	logger *zap.Logger
}

func NewProductService(repo repositories.ProductRepository, logger *zap.Logger) *ProductService {
	return &ProductService{repo: repo, logger: logger}
}

func (s *ProductService) GetCategoryNameByID(id int) (string, error) {
	return s.repo.GetCategoryNameByID(id)
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

	var categoryName string
	if product.CategoryID.Valid {
		categoryName, err = s.repo.GetCategoryNameByID(int(product.CategoryID.Int64))
		if err != nil {
			categoryName = ""
		}
	} else {
		categoryName = ""
	}

	response := utils.BuildProductResponse(product, categoryName)
	return &response, nil
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

	var categoryName string
	if updated.CategoryID.Valid {
		categoryName, err = s.repo.GetCategoryNameByIDTx(tx, int(updated.CategoryID.Int64))
		if err != nil {
			categoryName = ""
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	response := utils.BuildProductResponse(updated, categoryName)
	return &response, nil

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
		var categoryName string
		if p.CategoryID.Valid {
			categoryName, err = s.repo.GetCategoryNameByID(int(p.CategoryID.Int64))
			if err != nil {
				categoryName = ""
			}
		} else {
			categoryName = ""
		}

		response := utils.BuildProductResponse(&p, categoryName)
		result = append(result, response)

	}
	return result, nil
}

func (s *ProductService) GetFilteredRaw(
	search, category string,
	minPrice, maxPrice float64,
	limit, offset int,
	sort string,
	availability *bool,
) ([]models.Product, error) {
	return s.repo.GetFiltered(search, category, minPrice, maxPrice, limit, offset, sort, availability)
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

			var categoryName string
			if updated.CategoryID.Valid {
				categoryName, _ = s.repo.GetCategoryNameByIDTx(tx, int(updated.CategoryID.Int64))
			}

			response := utils.BuildProductResponse(updated, categoryName)
			responses = append(responses, response)

			continue
		}

		if err := s.repo.CreateTx(tx, &p); err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		var categoryName string
		if p.CategoryID.Valid {
			categoryName, _ = s.repo.GetCategoryNameByIDTx(tx, int(p.CategoryID.Int64))
		}

		response := utils.BuildProductResponse(&p, categoryName)
		responses = append(responses, response)

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

func (s *ProductService) PatchProduct(id int, updates map[string]interface{}) error {
	patch := buildProductPatch(updates)
	return s.repo.PatchProduct(id, patch)
}

func buildProductPatch(updates map[string]interface{}) models.ProductPatch {
	var patch models.ProductPatch

	if v, ok := updates["name"].(string); ok {
		patch.Name = &v
	}
	if v, ok := updates["description"].(string); ok {
		patch.Description = &v
	}
	if v, ok := updates["price"].(float64); ok {
		patch.Price = &v
	}
	if v, ok := updates["availability"].(bool); ok {
		patch.Availability = &v
	}
	if v, ok := updates["category_id"].(float64); ok {
		vi := int(v)
		patch.CategoryID = &vi
	}
	if v, ok := updates["url"].(string); ok {
		patch.Url = &v
	}

	return patch
}

func (s *ProductService) GetPaginated(
	search, category string,
	minPrice, maxPrice float64,
	limit, offset int,
	sort string,
	availability *bool,
) ([]models.ProductResponse, int, error) {
	products, err := s.GetFiltered(search, category, minPrice, maxPrice, limit, offset, sort, availability)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountFiltered(search, category, minPrice, maxPrice, availability)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) CountFiltered(
	search, category string,
	minPrice, maxPrice float64,
	availability *bool,
) (int, error) {
	return s.repo.CountFiltered(search, category, minPrice, maxPrice, availability)
}
