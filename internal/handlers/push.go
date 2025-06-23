package handlers

import (
	"chechnya-product/internal/models"
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
	DeleteSubscription(w http.ResponseWriter, r *http.Request)
	Subscribe(w http.ResponseWriter, r *http.Request)
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

// SendNotification
// @Summary Отправить push-уведомление
// @Description Отправляет уведомление одному пользователю по подписке
// @Tags Push
// @Accept json
// @Produce json
// @Param input body pushRequest true "Подписка и сообщение"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/push/send [post]
func (h *PushHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req pushRequest
	if req.Message == "" {
		req.Message = "🔔 Уведомление от Chechnya Product"
	}

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

// Broadcast
// @Summary Массовая рассылка push-уведомлений
// @Description Рассылает сообщение всем подписанным пользователям
// @Tags Push
// @Accept json
// @Produce json
// @Param input body map[string]string true "Сообщение для рассылки"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/push/broadcast [post]
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

// DeleteSubscription
// @Summary Удалить подписку
// @Description Удаляет подписку по endpoint
// @Tags Push
// @Produce json
// @Param endpoint query string true "URL подписки"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/push/delete [delete]
func (h *PushHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint == "" {
		h.logger.Warn("не указан endpoint для удаления")
		utils.ErrorJSON(w, http.StatusBadRequest, "Не указан endpoint")
		return
	}

	if err := h.service.DeleteByEndpoint(endpoint); err != nil {
		h.logger.Error("не удалось удалить подписку", zap.String("endpoint", endpoint), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка удаления подписки")
		return
	}

	h.logger.Info("подписка удалена", zap.String("endpoint", endpoint))
	utils.JSONResponse(w, http.StatusOK, "Подписка удалена", nil)
}

// Subscribe
// @Summary Подписка на push-уведомления
// @Description Регистрирует push-подписку пользователя, сохраняет её в базе
// @Tags Push
// @Accept json
// @Produce json
// @Param input body models.PushSubscriptionRequest true "Объект подписки"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/push/subscribe [post]
func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req models.PushSubscriptionRequest

	// Парсим JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Ошибка декодирования подписки", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	// Сохраняем подписку
	if err := h.service.SaveSubscription(req.Subscription, req.IsAdmin); err != nil {
		h.logger.Warn("Ошибка сохранения подписки", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Не удалось сохранить подписку")
		return
	}

	h.logger.Info("Подписка успешно сохранена", zap.String("endpoint", req.Subscription.Endpoint))
	utils.JSONResponse(w, http.StatusCreated, "Подписка успешно сохранена", nil)
}
