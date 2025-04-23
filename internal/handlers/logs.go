package handlers

import (
	"chechnya-product/internal/utils"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type LogHandlerInterface interface {
	GetLogs(w http.ResponseWriter, r *http.Request)
}

type LogHandler struct {
	logger  *zap.Logger
	logPath string
}

func NewLogHandler(logger *zap.Logger, logPath string) *LogHandler {
	return &LogHandler{logger: logger, logPath: logPath}
}

// GetLogs
// @Summary Получить лог-файл
// @Description Возвращает содержимое лог-файла (только для администратора)
// @Tags Логи
// @Security BearerAuth
// @Produce plain
// @Success 200 {string} string "Содержимое лог-файла"
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/logs [get]
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(h.logPath)
	if err != nil {
		h.logger.Error("failed to open log file", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Could not open log file")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
