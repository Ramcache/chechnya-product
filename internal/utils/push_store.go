package utils

import "chechnya-product/internal/models"

type PushStore interface {
	Save(sub models.SubscribeRequest)
	All() []models.SubscribeRequest
}

type InMemoryPushStore struct {
	data []models.SubscribeRequest
}

func NewInMemoryPushStore() *InMemoryPushStore {
	return &InMemoryPushStore{}
}

func (s *InMemoryPushStore) Save(sub models.SubscribeRequest) {
	s.data = append(s.data, sub)
}

func (s *InMemoryPushStore) All() []models.SubscribeRequest {
	return s.data
}
