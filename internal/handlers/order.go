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

// PlaceOrder godoc
// @Summary      Оформить заказ
// @Description  Создаёт заказ на основе содержимого корзины пользователя
// @Tags         orders
// @Security     BearerAuth
// @Produce      plain
// @Success      200 {string} string "Заказ успешно оформлен"
// @Failure      400 {string} string "Ошибка: корзина пуста или товара не хватает"
// @Router       /order [post]
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	err := h.service.PlaceOrder(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Заказ успешно оформлен"))
}

// GetUserOrders godoc
// @Summary      История заказов пользователя
// @Description  Возвращает список всех заказов текущего пользователя
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Success 200 {array} object
// @Failure      500 {string} string "Ошибка получения заказов"
// @Router       /orders [get]
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

// GetAllOrders godoc
// @Summary      Все заказы (только для админа)
// @Description  Возвращает список всех заказов в системе
// @Tags         admin-orders
// @Security     BearerAuth
// @Produce      json
// @Success 200 {array} object
// @Failure      500 {string} string "Ошибка получения заказов"
// @Router       /admin/orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Ошибка получения всех заказов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
