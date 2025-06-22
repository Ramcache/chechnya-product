package services

import (
	"bytes"
	"chechnya-product/config"
	"chechnya-product/internal/models"
	"chechnya-product/internal/repositories"
	"encoding/json"
	"errors"
	"github.com/SherClockHolmes/webpush-go"
	"go.uber.org/zap"
)

type PushServiceInterface interface {
	SendPush(sub webpush.Subscription, message string) error
	Broadcast(message string) error
	DeleteByEndpoint(endpoint string) error
}

type PushService struct {
	repo   repositories.PushRepositoryInterface
	logger *zap.Logger
	cfg    *config.Config
}

func NewPushService(repo repositories.PushRepositoryInterface, logger *zap.Logger, cfg *config.Config) *PushService {
	return &PushService{repo: repo, logger: logger, cfg: cfg}
}

func (s *PushService) SendPush(sub webpush.Subscription, message string) error {
	if message == "" {
		return errors.New("message is empty")
	}

	err := s.repo.SaveSubscription(models.Subscription{
		Endpoint: sub.Endpoint,
		P256dh:   sub.Keys.P256dh,
		Auth:     sub.Keys.Auth,
	})
	if err != nil {
		s.logger.Warn("не удалось сохранить подписку", zap.Error(err))
	}

	payload, _ := json.Marshal(map[string]string{
		"title": "Новое сообщение",
		"body":  message,
	})

	resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
		Subscriber:      "mailto:support@chechnya-product.ru",
		VAPIDPublicKey:  s.cfg.VAPIDPublicKey,
		VAPIDPrivateKey: s.cfg.VAPIDPrivateKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if resp.StatusCode >= 400 {
		s.logger.Error("webpush ошибка", zap.String("body", buf.String()))
		return errors.New("web push failed")
	}
	return nil
}

func (s *PushService) Broadcast(message string) error {
	subs, err := s.repo.GetAllSubscriptions()
	if err != nil {
		return err
	}
	for _, sub := range subs {
		webSub := webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}
		_ = s.SendPush(webSub, message)
	}
	return nil
}
func (s *PushService) DeleteByEndpoint(endpoint string) error {
	return s.repo.DeleteByEndpoint(endpoint)
}
