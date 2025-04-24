package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
)

type CategoryServiceInterface interface {
	GetAll() ([]models.Category, error)
	Create(name string, sortOrder int) error
	Update(id int, name string, sortOrder int) error
	Delete(id int) error
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
