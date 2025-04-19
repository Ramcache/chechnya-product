package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"chechnya-product/internal/repositories"
)

type VerificationService interface {
	StartVerification(ctx context.Context, phone string) (string, error)
	ConfirmCode(ctx context.Context, phone, code string) error
}

type verificationService struct {
	repo       repositories.VerificationRepository
	botNumber  string
	codeExpiry time.Duration
}

func NewVerificationService(repo repositories.VerificationRepository, botNumber string) VerificationService {
	return &verificationService{
		repo:       repo,
		botNumber:  botNumber,
		codeExpiry: 10 * time.Minute,
	}
}

func (s *verificationService) StartVerification(ctx context.Context, phone string) (string, error) {
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	if err := s.repo.SaveOrUpdate(ctx, phone, code); err != nil {
		return "", err
	}

	text := fmt.Sprintf("Мой код подтверждения: %s", code)
	encoded := url.QueryEscape(text)

	link := fmt.Sprintf("https://wa.me/%s?text=%s", s.botNumber, encoded)
	return link, nil
}

func (s *verificationService) ConfirmCode(ctx context.Context, phone, code string) error {
	v, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		return errors.New("verification not found")
	}

	if v.Code != code {
		return errors.New("invalid code")
	}

	if time.Since(v.CreatedAt) > s.codeExpiry {
		return errors.New("code expired")
	}

	return s.repo.Confirm(ctx, phone)
}
