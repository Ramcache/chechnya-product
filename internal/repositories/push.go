package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type PushRepositoryInterface interface {
	SaveSubscription(sub models.Subscription) error
	GetAllSubscriptions() ([]models.Subscription, error)
	DeleteByEndpoint(endpoint string) error
}

type PushRepository struct {
	db *sqlx.DB
}

func NewPushRepo(db *sqlx.DB) *PushRepository {
	return &PushRepository{db: db}
}

func (r *PushRepository) SaveSubscription(sub models.Subscription) error {
	_, err := r.db.Exec(`
	INSERT INTO push_subscriptions (endpoint, p256dh, auth, is_admin)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (endpoint)
	DO UPDATE SET p256dh = EXCLUDED.p256dh, auth = EXCLUDED.auth, is_admin = EXCLUDED.is_admin;
`, sub.Endpoint, sub.P256dh, sub.Auth, sub.IsAdmin)

	if err != nil {
		log.Println("❌ Ошибка при сохранении подписки:", err)
	}
	return err
}

func (r *PushRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	rows, err := r.db.Query(`SELECT endpoint, p256dh, auth, is_admin FROM push_subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		sub := models.Subscription{}
		if err := rows.Scan(&sub.Endpoint, &sub.P256dh, &sub.Auth, &sub.IsAdmin); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil

}

func (r *PushRepository) DeleteByEndpoint(endpoint string) error {
	_, err := r.db.Exec(`DELETE FROM push_subscriptions WHERE endpoint = $1`, endpoint)
	return err
}
