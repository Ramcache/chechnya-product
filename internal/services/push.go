package services

import (
	"chechnya-product/config"
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"log"
	"os"
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

func (s *PushService) SendPushToAdmins(order *models.Order) {
	// 1. Получить всех админов с подпиской
	adminSubs, err := s.repo.GetAdminSubscriptions()
	if err != nil {
		log.Printf("не удалось получить подписки админов: %v", err)
		return
	}

	// 2. Сформировать payload
	message := map[string]string{
		"title": "Новый заказ",
		"body":  fmt.Sprintf("Новый заказ #%d на сумму %.2f", order.ID, order.Total),
	}
	payload, _ := json.Marshal(message)

	// 3. Отправить всем админам
	for _, sub := range adminSubs {
		webpushSub := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256DH,
				Auth:   sub.Auth,
			},
		}

		resp, err := webpush.SendNotification(payload, webpushSub, &webpush.Options{
			TTL:             60,
			VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
			VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
			Subscriber:      "mailto:admin@example.com",
		})

		if err != nil {
			log.Printf("ошибка отправки пуша: %v", err)
			continue
		}
		defer resp.Body.Close()
	}
}
