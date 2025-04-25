package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/csv"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type OrderHandlerInterface interface {
	PlaceOrder(w http.ResponseWriter, r *http.Request)
	GetUserOrders(w http.ResponseWriter, r *http.Request)
	GetAllOrders(w http.ResponseWriter, r *http.Request)
	ExportOrdersCSV(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	service services.OrderServiceInterface
	logger  *zap.Logger
}

func NewOrderHandler(service services.OrderServiceInterface, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

// PlaceOrder
// @Summary Оформить заказ
// @Description Оформляет заказ из текущей корзины owner_id
// @Tags Заказ
// @Produce json
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/order [post]
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	if err := h.service.PlaceOrder(ownerID); err != nil {
		h.logger.Warn("failed to place order", zap.String("owner_id", ownerID), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Failed to place order")
		return
	}

	h.logger.Info("order placed", zap.String("owner_id", ownerID))
	utils.JSONResponse(w, http.StatusOK, "Order placed successfully", nil)
}

// GetUserOrders
// @Summary Получить заказы пользователя
// @Description Возвращает список заказов для текущего owner_id
// @Tags Заказ
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/orders [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	orders, err := h.service.GetOrders(ownerID)
	if err != nil {
		h.logger.Error("failed to get user orders", zap.String("owner_id", ownerID), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch user orders")
		return
	}

	h.logger.Info("user orders retrieved", zap.String("owner_id", ownerID), zap.Int("orders_count", len(orders)))
	utils.JSONResponse(w, http.StatusOK, "User orders retrieved", orders)
}

// GetAllOrders
// @Summary Получить все заказы (админ)
// @Description Возвращает список всех заказов (только для админа)
// @Tags Заказ
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		h.logger.Error("failed to get all orders", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch all orders")
		return
	}

	h.logger.Info("all orders retrieved", zap.Int("orders_count", len(orders)))
	utils.JSONResponse(w, http.StatusOK, "All orders retrieved", orders)
}

// ExportOrdersCSV
// @Summary Экспорт заказов в CSV (админ)
// @Description Экспортирует все заказы в формате CSV (только для админа)
// @Tags Заказ
// @Security BearerAuth
// @Produce text/csv
// @Success 200 {file} file "CSV файл"
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/orders/export [get]
func (h *OrderHandler) ExportOrdersCSV(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		h.logger.Error("failed to export orders to CSV", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch orders")
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
			utils.FormatFloat(order.Total),
			order.CreatedAt.Format(time.RFC3339),
		})
	}
}
