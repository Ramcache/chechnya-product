package handlers

import (
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type DashboardHandlerInterface interface {
	GetDashboard(w http.ResponseWriter, r *http.Request)
}

type DashboardHandler struct {
	service services.DashboardServiceInterface
	logger  *zap.Logger
}

func NewDashboardHandler(service services.DashboardServiceInterface, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{service: service, logger: logger}
}

// GetDashboard
// @Summary Дэшборд администратора
// @Description Возвращает метрики: заказы, выручка, топ товары, продажи по дням
// @Tags Админ / Дэшборд
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.DashboardData
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/dashboard [get]
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetDashboardData(r.Context())
	if err != nil {
		h.logger.Error("failed to load dashboard", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to load dashboard")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Dashboard data", data)
}
