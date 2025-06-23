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
	"regexp"
	"strings"
)

var base64URLRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

var base64urlPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type PushServiceInterface interface {
	SendPush(sub webpush.Subscription, message string) error
	Broadcast(message string) error
	DeleteByEndpoint(endpoint string) error
	SendPushToAdmins(message string) error
	SaveSubscription(sub webpush.Subscription, isAdmin bool) error
}

type PushService struct {
	repo   repositories.PushRepositoryInterface
	logger *zap.Logger
	cfg    *config.Config
}

func NewPushService(repo repositories.PushRepositoryInterface, logger *zap.Logger, cfg *config.Config) *PushService {
	return &PushService{repo: repo, logger: logger, cfg: cfg}
}

func (s *PushService) SaveSubscription(sub webpush.Subscription, isAdmin bool) error {
	// Проверка: ключи не пустые
	if sub.Keys.P256dh == "" || sub.Keys.Auth == "" {
		return errors.New("ключи подписки отсутствуют")
	}

	// Проверка: формат base64url
	if !base64urlPattern.MatchString(sub.Keys.P256dh) || !base64urlPattern.MatchString(sub.Keys.Auth) {
		s.logger.Warn("❌ Ключи не в формате base64url",
			zap.String("p256dh", sub.Keys.P256dh),
			zap.String("auth", sub.Keys.Auth),
		)
		return errors.New("ключи подписки имеют неверный формат (ожидается base64url)")
	}

	// Проверка длины
	if len(sub.Keys.P256dh) < 80 || len(sub.Keys.Auth) < 16 {
		return errors.New("ключи слишком короткие — возможно, подписка повреждена")
	}

	// Сохраняем подписку
	err := s.repo.SaveSubscription(models.Subscription{
		Endpoint: sub.Endpoint,
		P256dh:   sub.Keys.P256dh,
		Auth:     sub.Keys.Auth,
		IsAdmin:  isAdmin,
	})
	if err != nil {
		s.logger.Warn("❗ Не удалось сохранить подписку", zap.Error(err))
		return err
	}

	return nil
}

func (s *PushService) SendPush(sub webpush.Subscription, message string) error {
	// Подготавливаем payload
	payload, _ := json.Marshal(map[string]string{
		"title": "Новое сообщение",
		"body":  message,
	})

	s.logger.Debug("📦 Отправка пуша",
		zap.String("endpoint", sub.Endpoint),
		zap.Int("p256dh_len", len(sub.Keys.P256dh)),
		zap.Int("auth_len", len(sub.Keys.Auth)),
	)

	resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
		Subscriber:      "mailto:support@chechnya-product.ru",
		VAPIDPublicKey:  s.cfg.VAPIDPublicKey,
		VAPIDPrivateKey: s.cfg.VAPIDPrivateKey,
		TTL:             86400, // 1 день
	})

	if err != nil {
		s.logger.Error("❌ Webpush ошибка", zap.String("body", err.Error()))

		if strings.Contains(err.Error(), "unsubscribed") || strings.Contains(err.Error(), "expired") {
			_ = s.repo.DeleteByEndpoint(sub.Endpoint)
			s.logger.Info("🗑️ Удалена неактивная подписка", zap.String("endpoint", sub.Endpoint))
		}

		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		s.logger.Error("📛 Webpush ошибка",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", buf.String()),
		)
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

func (s *PushService) SendPushToAdmins(message string) error {
	subs, err := s.repo.GetAllSubscriptions()
	if err != nil {
		return err
	}

	adminCount := 0
	successCount := 0
	failCount := 0

	for _, sub := range subs {
		if !sub.IsAdmin {
			continue
		}
		adminCount++

		webSub := webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}

		err := s.SendPush(webSub, message)
		if err != nil {
			failCount++
			s.logger.Warn("❌ Ошибка отправки админу",
				zap.String("endpoint", sub.Endpoint),
				zap.Error(err))
		} else {
			successCount++
		}
	}

	s.logger.Info("📨 Push отправлен администраторам",
		zap.Int("admins", adminCount),
		zap.Int("успешно", successCount),
		zap.Int("ошибки", failCount),
	)
	return nil
}
