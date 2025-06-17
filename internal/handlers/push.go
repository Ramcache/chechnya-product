package handlers

import (
	"chechnya-product/internal/models"
	"chechnya-product/internal/utils"
	"encoding/json"
	"net/http"
)

type PushHandlerInterface interface {
	Subscribe(w http.ResponseWriter, r *http.Request)
}

type PushHandler struct {
	Store utils.PushStore // временно
}

func NewPushHandler(store utils.PushStore) *PushHandler {
	return &PushHandler{Store: store}
}

func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var sub models.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Неверный формат подписки")
		return
	}

	// временно просто сохраняем в памяти
	h.Store.Save(sub)

	utils.JSONResponse(w, http.StatusCreated, "Подписка сохранена", nil)
}
