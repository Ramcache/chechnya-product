package services

import (
	"bytes"
	"chechnya-product/internal/repositories"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type VerificationService struct {
	repo     repositories.VerificationRepository
	whatsapp string // номер бота
}

func NewVerificationService(repo repositories.VerificationRepository, whatsappBotPhone string) VerificationService {
	return VerificationService{repo: repo, whatsapp: whatsappBotPhone}
}

func (s VerificationService) GenerateCode(phone string) string {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	_ = s.repo.SaveCode(phone, code, 5*time.Minute)
	return code
}

func (s VerificationService) VerifyCode(phone, code string) error {
	stored, err := s.repo.GetCode(phone)
	if err != nil {
		return fmt.Errorf("code not found or expired")
	}
	if stored != code {
		return fmt.Errorf("wrong code")
	}
	_ = s.repo.DeleteCode(phone)
	return nil
}

func (s VerificationService) SendCodeViaWhatsApp(phone, code string) {
	body := map[string]string{"phone": phone, "code": code}
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post("http://localhost:5555/send-code", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("Failed to send WhatsApp message:", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func (s VerificationService) MarkPhoneVerified(phone string) error {
	return s.repo.MarkVerified(phone)
}
