package services

import "chechnya-product/internal/repositories"

type AdminServiceInterface interface {
	TruncateTable(tableName string) error
}

type AdminService struct {
	repo repositories.AdminRepoInterface
}

func NewAdminService(repo repositories.AdminRepoInterface) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) TruncateTable(tableName string) error {
	return s.repo.TruncateTable(tableName)
}
