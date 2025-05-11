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
	TruncateAllTablesHandler(w http.ResponseWriter, r *http.Request)
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

// TruncateAllTablesHandler
// @Summary Очистить все таблицы (только админ)
// @Security BearerAuth
// @Tags Admin
// @Produce json
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/admin/truncate/all [post]
func (h *AdminHandler) TruncateAllTablesHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.service.TruncateAllTables(); err != nil {
		h.logger.Error("ошибка очистки всех таблиц", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Ошибка очистки всех таблиц: "+err.Error())
		return
	}

	h.logger.Info("Все таблицы очищены")
	utils.JSONResponse(w, http.StatusOK, "Все таблицы очищены", nil)
}
