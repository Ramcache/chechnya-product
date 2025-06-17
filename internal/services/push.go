package services

import (
	"chechnya-product/config"
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"encoding/json"
	"errors"
	"github.com/SherClockHolmes/webpush-go"
)

type PushServiceInterface interface {
	Save(sub *models.PushSubscription) error
	SendToAll(title, message string) error
}

type PushService struct {
	repo   repositories.PushRepository
	config *config.Config
}

func NewPushService(repo repositories.PushRepository, cfg *config.Config) *PushService {
	return &PushService{repo: repo, config: cfg}
}

func (s *PushService) Save(sub *models.PushSubscription) error {
	return s.repo.Save(sub)
}

func (s *PushService) SendToAll(title, message string) error {
	subs, err := s.repo.GetAllSubscriptions()
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		return errors.New("нет подписчиков")
	}
	vapidPublicKey := s.config.VAPIDPublicKey
	vapidPrivateKey := s.config.VAPIDPrivateKey
	vapidEmail := "mailto:ramaro@internet.ru"

	payload, _ := json.Marshal(map[string]string{
		"title":   title,
		"message": message,
	})

	for _, sub := range subs {
		subscription := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256DH,
				Auth:   sub.Auth,
			},
		}

		resp, err := webpush.SendNotification(
			payload,
			subscription,
			&webpush.Options{
				Subscriber:      vapidEmail,
				VAPIDPublicKey:  vapidPublicKey,
				VAPIDPrivateKey: vapidPrivateKey,
				TTL:             30,
			},
		)
		if err != nil {
			// Можно добавить лог: удалить подписку, если она устарела
			continue
		}
		resp.Body.Close()
	}

	return nil
}
