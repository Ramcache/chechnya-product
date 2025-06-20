package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/models"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type PushHandlerInterface interface {
	SaveSubscription(w http.ResponseWriter, r *http.Request)
	SendPush(w http.ResponseWriter, r *http.Request)
}

type PushHandler struct {
	service services.PushServiceInterface
	logger  *zap.Logger
}

func NewPushHandler(service services.PushServiceInterface, logger *zap.Logger) *PushHandler {
	return &PushHandler{service: service, logger: logger}
}

type SaveSubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
	P256DH   string `json:"p256dh"`
	Auth     string `json:"auth"`
}

type SendPushRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SaveSubscription
// @Summary Сохранить push-подписку
// @Tags Push
// @Accept json
// @Produce json
// @Param input body SaveSubscriptionRequest true "Подписка"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/push/subscribe [post]
func (h *PushHandler) SaveSubscription(w http.ResponseWriter, r *http.Request) {
	var req SaveSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	claims := middleware.GetUserClaims(r)
	if claims == nil {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Не авторизован")
		return
	}
	userID := claims.UserID

	sub := &models.PushSubscription{
		Endpoint: req.Endpoint,
		P256DH:   req.P256DH,
		Auth:     req.Auth,
		UserID:   userID,
	}

	if err := h.service.Save(sub); err != nil {
		h.logger.Error("не удалось сохранить подписку", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка сохранения")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Подписка сохранена", nil)
}

// SendPush
// @Summary Отправить push-уведомление
// @Tags Push
// @Accept json
// @Produce json
// @Param input body SendPushRequest true "Уведомление"
// @Success 200 {object} utils.SuccessResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/push/send [post]
func (h *PushHandler) SendPush(w http.ResponseWriter, r *http.Request) {
	var req SendPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Неверный формат")
		return
	}

	if err := h.service.SendToAll(req.Title, req.Body); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка отправки")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Уведомления отправлены", nil)
}
