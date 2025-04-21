package handlers

import (
	"chechnya-product/internal/services"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type VerificationHandler struct {
	service services.VerificationService
	logger  *zap.Logger
}

func NewVerificationHandler(service services.VerificationService, logger *zap.Logger) *VerificationHandler {
	return &VerificationHandler{service: service, logger: logger}
}

type StartRequest struct {
	Phone string `json:"phone"`
}

type ConfirmRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// POST /verify/start
func (h *VerificationHandler) StartVerification(w http.ResponseWriter, r *http.Request) {
	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	code := h.service.GenerateCode(req.Phone)
	go h.service.SendCodeViaWhatsApp(req.Phone, code)

	h.logger.Info("Verification code sent", zap.String("phone", req.Phone), zap.String("code", code))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verification code sent via WhatsApp"))
}

// POST /verify/confirm
func (h *VerificationHandler) ConfirmCode(w http.ResponseWriter, r *http.Request) {
	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.VerifyCode(req.Phone, req.Code); err != nil {
		h.logger.Warn("code verification failed", zap.String("phone", req.Phone), zap.String("code", req.Code))
		http.Error(w, "Invalid or expired code", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkPhoneVerified(req.Phone); err != nil {
		h.logger.Error("failed to mark user verified", zap.Error(err))
		http.Error(w, "Failed to verify user", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Phone number verified", zap.String("phone", req.Phone))
	w.Write([]byte("Phone verified successfully"))
}
