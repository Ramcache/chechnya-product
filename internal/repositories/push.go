package repositories

import (
	"chechnya-product/internal/models"
	"github.com/jmoiron/sqlx"
)

type PushRepository interface {
	Save(sub *models.PushSubscription) error
	GetAllByUser(userID int) ([]models.PushSubscription, error)
	GetAllSubscriptions() ([]models.PushSubscription, error)
	GetAdminSubscriptions() ([]models.PushSubscription, error)
}

type pushRepo struct {
	db *sqlx.DB
}

func NewPushRepo(db *sqlx.DB) PushRepository {
	return &pushRepo{db: db}
}

func (r *pushRepo) Save(sub *models.PushSubscription) error {
	query := `
	INSERT INTO push_subscriptions (endpoint, p256dh, auth, user_id)
	VALUES (:endpoint, :p256dh, :auth, :user_id)
	`
	_, err := r.db.NamedExec(query, sub)
	return err
}

func (r *pushRepo) GetAllByUser(userID int) ([]models.PushSubscription, error) {
	var subs []models.PushSubscription
	err := r.db.Select(&subs, "SELECT * FROM push_subscriptions WHERE user_id=$1", userID)
	return subs, err
}

func (r *pushRepo) GetAllSubscriptions() ([]models.PushSubscription, error) {
	var subs []models.PushSubscription
	err := r.db.Select(&subs, `SELECT endpoint, p256dh, auth FROM push_subscriptions`)
	return subs, err
}

func (r *pushRepo) GetAdminSubscriptions() ([]models.PushSubscription, error) {
	var subs []models.PushSubscription
	err := r.db.Select(&subs, `
		SELECT ps.*
		FROM push_subscriptions ps
		JOIN users u ON ps.user_id = u.id
		WHERE u.role = 'admin'
	`)
	return subs, err
}
