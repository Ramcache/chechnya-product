package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	err := h.service.PlaceOrder(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Заказ успешно оформлен"))
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	orders, err := h.service.GetOrders(userID)
	if err != nil {
		http.Error(w, "Ошибка получения заказов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Ошибка получения всех заказов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
