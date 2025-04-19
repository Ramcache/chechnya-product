package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"

	"chechnya-product/internal/models"
)

type VerificationRepository interface {
	SaveOrUpdate(ctx context.Context, phone, code string) error
	GetByPhone(ctx context.Context, phone string) (*models.Verification, error)
	Confirm(ctx context.Context, phone string) error
}

type verificationRepo struct {
	db *sqlx.DB
}

func NewVerificationRepository(db *sqlx.DB) VerificationRepository {
	return &verificationRepo{db: db}
}

func (r *verificationRepo) SaveOrUpdate(ctx context.Context, phone, code string) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO phone_verifications (phone, code, created_at, confirmed)
        VALUES ($1, $2, NOW(), FALSE)
        ON CONFLICT (phone)
        DO UPDATE SET code = EXCLUDED.code, created_at = EXCLUDED.created_at, confirmed = FALSE
    `, phone, code)
	return err
}

func (r *verificationRepo) GetByPhone(ctx context.Context, phone string) (*models.Verification, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT phone, code, created_at, confirmed
        FROM phone_verifications
        WHERE phone = $1
    `, phone)

	var v models.Verification
	if err := row.Scan(&v.Phone, &v.Code, &v.CreatedAt, &v.Confirmed); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *verificationRepo) Confirm(ctx context.Context, phone string) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE phone_verifications SET confirmed = TRUE WHERE phone = $1
    `, phone)
	return err
}
