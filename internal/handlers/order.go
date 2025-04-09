package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

// ExportOrdersCSV godoc
// @Summary      Экспорт заказов в CSV
// @Description  Возвращает CSV-файл со всеми заказами (только для админа)
// @Tags         admin-orders
// @Security     BearerAuth
// @Produce      text/csv
// @Success      200
// @Failure      500 {string} string "Ошибка экспорта"
// @Router       /admin/orders/export [get]
func (h *OrderHandler) ExportOrdersCSV(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Ошибка получения заказов", http.StatusInternalServerError)
		return
	}

	// Заголовки для скачивания
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=orders.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Заголовок CSV
	writer.Write([]string{"Order ID", "User ID", "Total", "Created At"})

	for _, order := range orders {
		writer.Write([]string{
			strconv.Itoa(order.ID),
			strconv.Itoa(order.UserID),
			fmt.Sprintf("%.2f", order.Total),
			order.CreatedAt.Format(time.RFC3339),
		})
	}
}
