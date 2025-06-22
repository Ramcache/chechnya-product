package handlers

import (
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"github.com/SherClockHolmes/webpush-go"
	"go.uber.org/zap"
	"net/http"
)

type PushHandlerInterface interface {
	SendNotification(w http.ResponseWriter, r *http.Request)
	Broadcast(w http.ResponseWriter, r *http.Request)
}

type PushHandler struct {
	service services.PushServiceInterface
	logger  *zap.Logger
}

func NewPushHandler(service services.PushServiceInterface, logger *zap.Logger) *PushHandler {
	return &PushHandler{service: service, logger: logger}
}

type pushRequest struct {
	Subscription webpush.Subscription `json:"subscription"`
	Message      string               `json:"message"`
}

func (h *PushHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req pushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("невалидный JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	if err := h.service.SendPush(req.Subscription, req.Message); err != nil {
		h.logger.Error("ошибка отправки push", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка отправки push")
		return
	}

	h.logger.Info("push отправлен")
	utils.JSONResponse(w, http.StatusOK, "Push отправлен", nil)
}

func (h *PushHandler) Broadcast(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("невалидный JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	if err := h.service.Broadcast(req.Message); err != nil {
		h.logger.Error("ошибка рассылки", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка рассылки push")
		return
	}

	h.logger.Info("рассылка завершена")
	utils.JSONResponse(w, http.StatusOK, "Рассылка завершена", nil)
}
