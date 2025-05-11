package handlers

import (
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type AdminInterface interface {
	TruncateTableHandler(w http.ResponseWriter, r *http.Request)
}

type AdminHandler struct {
	service services.AdminServiceInterface
	logger  *zap.Logger
}

func NewAdminHandler(service services.AdminServiceInterface, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{service: service, logger: logger}
}

type TruncateRequest struct {
	Table string `json:"table"`
}

// TruncateTableHandler
// @Summary Очистка таблицы с перезапуском ID
// @Security BearerAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param table body TruncateRequest true "Название таблицы для очистки"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/admin/truncate [post]
func (h *AdminHandler) TruncateTableHandler(w http.ResponseWriter, r *http.Request) {
	var req TruncateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Ошибка декодирования")
		return
	}

	if err := h.service.TruncateTable(req.Table); err != nil {
		h.logger.Error("ошибка очистки таблицы", zap.String("table", req.Table), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Ошибка: "+err.Error())
		return
	}

	h.logger.Info("Таблица очищена", zap.String("table", req.Table))
	utils.JSONResponse(w, http.StatusOK, "Таблица очищена", nil)
}
