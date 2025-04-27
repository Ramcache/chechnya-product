package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/utils"
	"fmt"
	"go.uber.org/zap"
)

type CategoryServiceInterface interface {
	GetAll() ([]models.Category, error)
	Create(name string, sortOrder int) (*models.Category, error)
	Update(id int, name string, sortOrder int) error
	Delete(id int) error
	CreateBulk(categories []utils.CategoryRequest) ([]models.Category, error)
	PartialUpdate(id int, name *string, sortOrder *int) (*models.Category, error)
}

type CategoryService struct {
	repo   repositories.CategoryRepository
	logger *zap.Logger
}

func NewCategoryService(repo repositories.CategoryRepository, logger *zap.Logger) *CategoryService {
	return &CategoryService{repo: repo, logger: logger}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(name string, sortOrder int) (*models.Category, error) {
	category := &models.Category{
		Name:      name,
		SortOrder: sortOrder,
	}
	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(id int, name string, sortOrder int) error {
	return s.repo.Update(id, name, sortOrder)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) CreateBulk(categories []utils.CategoryRequest) ([]models.Category, error) {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	var created []models.Category
	var txErr error

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if txErr != nil {
			_ = tx.Rollback()
		}
	}()

	for _, cat := range categories {
		if cat.Name == "" {
			txErr = fmt.Errorf("category name cannot be empty")
			return nil, txErr
		}

		existing, err := s.repo.GetByNameTx(tx, cat.Name)
		if err == nil && existing != nil {
			s.logger.Info("category already exists, skipping", zap.String("name", cat.Name))
			continue
		}

		newCat, err := s.repo.CreateReturningTx(tx, cat.Name, cat.SortOrder)
		if err != nil {
			txErr = fmt.Errorf("failed to create category '%s': %w", cat.Name, err)
			return nil, txErr
		}

		s.logger.Info("category created", zap.String("name", newCat.Name), zap.Int("id", newCat.ID))
		created = append(created, *newCat)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return created, nil
}

func (s *CategoryService) PartialUpdate(id int, name *string, sortOrder *int) (*models.Category, error) {
	err := s.repo.PartialUpdate(id, name, sortOrder)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
