package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"context"
)

type DashboardServiceInterface interface {
	GetDashboardData(ctx context.Context) (*models.DashboardData, error)
}

type DashboardService struct {
	repo repositories.DashboardRepositoryInterface
}

func NewDashboardService(repo repositories.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetDashboardData(ctx context.Context) (*models.DashboardData, error) {
	return s.repo.GetDashboardData(ctx)
}
