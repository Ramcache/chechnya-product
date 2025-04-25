package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/utils"
	"fmt"
)

type CategoryServiceInterface interface {
	GetAll() ([]models.Category, error)
	Create(name string, sortOrder int) error
	Update(id int, name string, sortOrder int) error
	Delete(id int) error
	CreateBulk(categories []utils.CategoryRequest) error
}

type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
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

func (s *CategoryService) CreateBulk(categories []utils.CategoryRequest) error {
	for _, cat := range categories {
		if cat.Name == "" {
			return fmt.Errorf("category name cannot be empty")
		}
		err := s.repo.Create(cat.Name, cat.SortOrder)
		if err != nil {
			return err
		}
	}
	return nil
}
