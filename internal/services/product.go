package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
)

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) AddProduct(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) UpdateProduct(id int, product *models.Product) error {
	return s.repo.Update(id, product)
}
