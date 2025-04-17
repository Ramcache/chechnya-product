package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/csv"
	"encoding/json"
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

// PlaceOrder
// @Summary Оформить заказ
// @Description Создаёт новый заказ из товаров в корзине пользователя
// @Security BearerAuth
// @Produce plain
// @Success 200 {string} string "Заказ успешно создан"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/order [post]
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	if err := h.service.PlaceOrder(userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Order placed successfully"))
}

// GetUserOrders
// @Summary Получить заказы пользователя
// @Description Возвращает список заказов текущего пользователя
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/orders [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	orders, err := h.service.GetOrders(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user orders", http.StatusInternalServerError)
		return
	}

	writeJSON(w, orders)
}

// GetAllOrders
// @Summary Получить все заказы (Админ)
// @Description Возвращает список всех заказов (только для администратора)
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to fetch all orders", http.StatusInternalServerError)
		return
	}

	writeJSON(w, orders)
}

// ExportOrdersCSV
// @Summary Экспортировать заказы в CSV (Админ)
// @Description Экспортирует все заказы в формате CSV (только для администратора)
// @Security BearerAuth
// @Produce text/csv
// @Success 200 {file} file "CSV-файл с заказами"
// @Failure 500 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/admin/orders/export [get]
func (h *OrderHandler) ExportOrdersCSV(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=orders.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Order ID", "User ID", "Total", "Created At"})

	for _, order := range orders {
		writer.Write([]string{
			strconv.Itoa(order.ID),
			strconv.Itoa(order.UserID),
			formatFloat(order.Total),
			order.CreatedAt.Format(time.RFC3339),
		})
	}
}

// writeJSON - универсальный хелпер для JSON-ответов
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// formatFloat - вывод float с двумя знаками после запятой
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
