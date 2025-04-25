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
	Create(name string, sortOrder int) error
	Update(id int, name string, sortOrder int) error
	Delete(id int) error
	CreateBulk(categories []utils.CategoryRequest) ([]models.Category, error)
	PartialUpdate(id int, name *string, sortOrder *int) error
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

func (s *CategoryService) Create(name string, sortOrder int) error {
	return s.repo.Create(name, sortOrder)
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
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	var created []models.Category

	for _, cat := range categories {
		if cat.Name == "" {
			_ = tx.Rollback()
			return nil, fmt.Errorf("category name cannot be empty")
		}

		existing, err := s.repo.GetByNameTx(tx, cat.Name)
		if err == nil && existing != nil {
			s.logger.Info("category already exists, skipping", zap.String("name", cat.Name))
			continue
		}

		newCat, err := s.repo.CreateReturningTx(tx, cat.Name, cat.SortOrder)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("failed to create category '%s': %w", cat.Name, err)
		}

		s.logger.Info("category created", zap.String("name", newCat.Name), zap.Int("id", newCat.ID))
		created = append(created, *newCat)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return created, nil
}

func (s *CategoryService) PartialUpdate(id int, name *string, sortOrder *int) error {
	return s.repo.PartialUpdate(id, name, sortOrder)
}
