package services

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/ws"
)

type AnnouncementServiceInterface interface {
	GetAll() ([]models.Announcement, error)
	GetByID(id int) (*models.Announcement, error)
	Create(title, content string) (*models.Announcement, error)
	Update(id int, title, content string) error
	Delete(id int) error
}

type AnnouncementService struct {
	repo repositories.AnnouncementRepository
	hub  *ws.Hub // для WebSocket рассылки
}

func NewAnnouncementService(repo repositories.AnnouncementRepository, hub *ws.Hub) *AnnouncementService {
	return &AnnouncementService{repo: repo, hub: hub}
}

func (s *AnnouncementService) GetByID(id int) (*models.Announcement, error) {
	return s.repo.GetByID(id)
}

func (s *AnnouncementService) Create(title, content string) (*models.Announcement, error) {
	ann, err := s.repo.Create(title, content)
	if err != nil {
		return nil, err
	}
	s.hub.BroadcastAnnouncement(*ann)
	return ann, nil
}

func (s *AnnouncementService) GetAll() ([]models.Announcement, error) {
	return s.repo.GetAll()
}

func (s *AnnouncementService) Update(id int, title, content string) error {
	return s.repo.Update(id, title, content)
}

func (s *AnnouncementService) Delete(id int) error {
	return s.repo.Delete(id)
}
