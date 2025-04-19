package handlers

import (
	"encoding/json"
	"net/http"

	"chechnya-product/internal/services"
	"go.uber.org/zap"
)

// VerificationHandler handles phone verification logic
type VerificationHandler struct {
	service services.VerificationService
	log     *zap.Logger
}

// NewVerificationHandler creates a new instance of VerificationHandler
func NewVerificationHandler(service services.VerificationService, log *zap.Logger) *VerificationHandler {
	return &VerificationHandler{
		service: service,
		log:     log,
	}
}

// StartRequest represents a request to start phone verification
type StartRequest struct {
	Phone string `json:"phone" example:"+79991234567"`
}

// StartResponse contains a WhatsApp link for verification
type StartResponse struct {
	WhatsAppLink string `json:"whatsapp_link" example:"https://wa.me/79991234567?text=Мой%20код%20подтверждения%3A%20123456"`
}

// ConfirmRequest represents a request to confirm a phone verification
type ConfirmRequest struct {
	Phone string `json:"phone" example:"+79991234567"`
	Code  string `json:"code" example:"123456"`
}

// StartVerification godoc
//
//	@Summary		Начать подтверждение телефона
//	@Description	Генерирует код и возвращает ссылку на WhatsApp с готовым сообщением
//	@Tags			verification
//	@Accept			json
//	@Produce		json
//	@Param			request	body		StartRequest	true	"Телефон пользователя"
//	@Success		200		{object}	StartResponse
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		500		{string}	string	"Failed to start verification"
//	@Router			/verify/start [post]
func (h *VerificationHandler) StartVerification(w http.ResponseWriter, r *http.Request) {
	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	link, err := h.service.StartVerification(r.Context(), req.Phone)
	if err != nil {
		h.log.Error("Failed to start verification", zap.String("phone", req.Phone), zap.Error(err))
		http.Error(w, "Failed to start verification", http.StatusInternalServerError)
		return
	}

	h.log.Info("Verification code generated", zap.String("phone", req.Phone))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StartResponse{WhatsAppLink: link})
}

// ConfirmCode godoc
//
//	@Summary		Подтвердить код
//	@Description	Проверяет код подтверждения, подтверждает номер телефона
//	@Tags			verification
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ConfirmRequest	true	"Телефон и код"
//	@Success		200		{string}	string	"Phone confirmed"
//	@Failure		400		{string}	string	"Invalid request or code"
//	@Router			/verify/confirm [post]
func (h *VerificationHandler) ConfirmCode(w http.ResponseWriter, r *http.Request) {
	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.ConfirmCode(r.Context(), req.Phone, req.Code); err != nil {
		h.log.Warn("Verification failed", zap.String("phone", req.Phone), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Info("Phone successfully confirmed", zap.String("phone", req.Phone))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Phone confirmed"))
}
