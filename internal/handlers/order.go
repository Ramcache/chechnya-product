package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/csv"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type OrderHandler struct {
	service *services.OrderService
	logger  *zap.Logger
}

func NewOrderHandler(service *services.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

// PlaceOrder
// @Summary Оформить заказ
// @Description Оформляет заказ из текущей корзины owner_id
// @Tags Заказ
// @Produce plain
// @Success 200 {string} string "Order placed successfully"
// @Failure 400 {object} ErrorResponse
// @Router /api/order [post]
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	if err := h.service.PlaceOrder(ownerID); err != nil {
		h.logger.Warn("failed to place order", zap.String("owner_id", ownerID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("order placed", zap.String("owner_id", ownerID))
	w.Write([]byte("Order placed successfully"))
}

// GetUserOrders
// @Summary Получить заказы пользователя
// @Description Возвращает список заказов для текущего owner_id
// @Tags Заказ
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse
// @Router /api/orders [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	orders, err := h.service.GetOrders(ownerID)
	if err != nil {
		h.logger.Error("failed to get user orders", zap.String("owner_id", ownerID), zap.Error(err))
		http.Error(w, "Failed to fetch user orders", http.StatusInternalServerError)
		return
	}

	h.logger.Info("user orders retrieved", zap.String("owner_id", ownerID), zap.Int("orders_count", len(orders)))
	writeJSON(w, orders)
}

// GetAllOrders
// @Summary Получить все заказы (админ)
// @Description Возвращает список всех заказов (только для админа)
// @Tags Админ / Заказы
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		h.logger.Error("failed to get all orders", zap.Error(err))
		http.Error(w, "Failed to fetch all orders", http.StatusInternalServerError)
		return
	}

	h.logger.Info("all orders retrieved", zap.Int("orders_count", len(orders)))
	writeJSON(w, orders)
}

// ExportOrdersCSV
// @Summary Экспорт заказов в CSV (админ)
// @Description Экспортирует все заказы в формате CSV (только для админа)
// @Tags Админ / Заказы
// @Security BearerAuth
// @Produce text/csv
// @Success 200 {string} string "CSV файл"
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/orders/export [get]
func (h *OrderHandler) ExportOrdersCSV(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		h.logger.Error("failed to export orders to CSV", zap.Error(err))
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	h.logger.Info("exporting orders to CSV", zap.Int("orders_count", len(orders)))

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=orders.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Order ID", "Owner ID", "Total", "Created At"})

	for _, order := range orders {
		writer.Write([]string{
			strconv.Itoa(order.ID),
			order.OwnerID,
			formatFloat(order.Total),
			order.CreatedAt.Format(time.RFC3339),
		})
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
