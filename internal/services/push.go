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
)

var base64URLRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

var base64urlPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type PushServiceInterface interface {
	SendPush(sub webpush.Subscription, message string, isAdmin bool) error
	Broadcast(message string) error
	DeleteByEndpoint(endpoint string) error
	SendPushToAdmins(message string) error
}

type PushService struct {
	repo   repositories.PushRepositoryInterface
	logger *zap.Logger
	cfg    *config.Config
}

func NewPushService(repo repositories.PushRepositoryInterface, logger *zap.Logger, cfg *config.Config) *PushService {
	return &PushService{repo: repo, logger: logger, cfg: cfg}
}

func (s *PushService) SendPush(sub webpush.Subscription, message string, isAdmin bool) error {
	if message == "" {
		return errors.New("message is empty")
	}
	if !base64urlPattern.MatchString(sub.Keys.P256dh) || !base64urlPattern.MatchString(sub.Keys.Auth) {
		s.logger.Warn("❌ Ключи не в формате base64url",
			zap.String("p256dh", sub.Keys.P256dh),
			zap.String("auth", sub.Keys.Auth),
		)
		return errors.New("ключи подписки имеют неверный формат (ожидается base64url)")
	}
	err := s.repo.SaveSubscription(models.Subscription{
		Endpoint: sub.Endpoint,
		P256dh:   sub.Keys.P256dh,
		Auth:     sub.Keys.Auth,
		IsAdmin:  isAdmin,
	})

	if err != nil {
		s.logger.Warn("не удалось сохранить подписку", zap.Error(err))
	}

	payload, _ := json.Marshal(map[string]string{
		"title": "Новое сообщение",
		"body":  message,
	})
	if sub.Keys.P256dh == "" || sub.Keys.Auth == "" {
		return errors.New("ключи подписки отсутствуют")
	}

	if !base64URLRegex.MatchString(sub.Keys.P256dh) || !base64URLRegex.MatchString(sub.Keys.Auth) {
		s.logger.Warn("невалидный формат ключей", zap.String("p256dh", sub.Keys.P256dh), zap.String("auth", sub.Keys.Auth))
		return errors.New("невалидный формат ключей (ожидается base64url)")
	}

	if len(sub.Keys.P256dh) < 80 || len(sub.Keys.Auth) < 16 {
		return errors.New("ключи слишком короткие — возможно, подписка повреждена")
	}
	s.logger.Debug("📦 Входящая подписка",
		zap.String("endpoint", sub.Endpoint),
		zap.String("p256dh", sub.Keys.P256dh),
		zap.String("auth", sub.Keys.Auth),
		zap.Int("p256dh_len", len(sub.Keys.P256dh)),
		zap.Int("auth_len", len(sub.Keys.Auth)),
	)

	resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
		Subscriber:      "mailto:support@chechnya-product.ru",
		VAPIDPublicKey:  s.cfg.VAPIDPublicKey,
		VAPIDPrivateKey: s.cfg.VAPIDPrivateKey,
		TTL:             86400,
	})
	s.logger.Debug("🚀 Отправка пуша через webpush", zap.String("endpoint", sub.Endpoint))

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
		_ = s.SendPush(webSub, message, true)
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
		_ = s.SendPush(webSub, message, true) // 3-й аргумент можно не важен, т.к. это отправка, не сохранение
	}

	s.logger.Info("📨 Push отправлен администраторам", zap.Int("admin_count", adminCount))
	return nil
}
