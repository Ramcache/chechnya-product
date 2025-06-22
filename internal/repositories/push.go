package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type PushRepositoryInterface interface {
	SaveSubscription(sub models.Subscription) error
	GetAllSubscriptions() ([]models.Subscription, error)
}

type PushRepository struct {
	db *sqlx.DB
}

func NewPushRepo(db *sqlx.DB) *PushRepository {
	return &PushRepository{db: db}
}

func (r *PushRepository) SaveSubscription(sub models.Subscription) error {
	_, err := r.db.Exec(`
	INSERT INTO push_subscriptions (endpoint, p256dh, auth)
	VALUES ($1, $2, $3)
	ON CONFLICT (endpoint)
	DO UPDATE SET p256dh = EXCLUDED.p256dh, auth = EXCLUDED.auth;
	`, sub.Endpoint, sub.P256dh, sub.Auth)

	if err != nil {
		log.Println("❌ Ошибка при сохранении подписки:", err)
	}
	return err
}

func (r *PushRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	rows, err := r.db.Query(`SELECT endpoint, p256dh, auth FROM push_subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		sub := models.Subscription{}
		if err := rows.Scan(&sub.Endpoint, &sub.P256dh, &sub.Auth); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
