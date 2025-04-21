package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type VerificationRepository interface {
	SaveCode(phone, code string, ttl time.Duration) error
	GetCode(phone string) (string, error)
	DeleteCode(phone string) error
	MarkVerified(phone string) error
}

type verificationRepo struct {
	db *sqlx.DB
}

func NewVerificationRepository(db *sqlx.DB) *verificationRepo {
	return &verificationRepo{db: db}
}

func (r *verificationRepo) SaveCode(phone, code string, ttl time.Duration) error {
	expiry := time.Now().Add(ttl)
	_, err := r.db.Exec(`
		INSERT INTO verification_codes (phone, code, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (phone)
		DO UPDATE SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at
	`, phone, code, expiry)
	return err
}

func (r *verificationRepo) GetCode(phone string) (string, error) {
	var (
		code      string
		expiresAt time.Time
	)
	err := r.db.QueryRow(`
		SELECT code, expires_at FROM verification_codes WHERE phone = $1
	`, phone).Scan(&code, &expiresAt)

	if err != nil {
		return "", err
	}
	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("code expired")
	}
	return code, nil
}

func (r *verificationRepo) DeleteCode(phone string) error {
	_, err := r.db.Exec(`DELETE FROM verification_codes WHERE phone = $1`, phone)
	return err
}

func (r *verificationRepo) MarkVerified(phone string) error {
	_, err := r.db.Exec(`UPDATE users SET is_verified = TRUE WHERE phone = $1`, phone)
	return err
}
