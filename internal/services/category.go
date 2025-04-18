package services

import "chechnya-product/internal/repositories"

type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]string, error) {
	categories, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, c := range categories {
		names = append(names, c.Name)
	}
	return names, nil
}

func (s *CategoryService) Create(name string) error {
	return s.repo.Create(name)
}

func (s *CategoryService) Update(id int, name string) error {
	return s.repo.Update(id, name)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
