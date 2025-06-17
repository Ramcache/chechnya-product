package utils

import "chechnya-product/internal/models"

type PushStore interface {
	Save(sub models.PushSubscription)
	All() []models.PushSubscription
}

type InMemoryPushStore struct {
	data []models.PushSubscription
}

func NewInMemoryPushStore() *InMemoryPushStore {
	return &InMemoryPushStore{}
}

func (s *InMemoryPushStore) Save(sub models.PushSubscription) {
	s.data = append(s.data, sub)
}

func (s *InMemoryPushStore) All() []models.PushSubscription {
	return s.data
}
