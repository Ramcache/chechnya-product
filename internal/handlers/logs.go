package handlers

import (
	"chechnya-product/internal/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
// @Summary      Получить лог-файл
// @Description  Возвращает лог за указанную дату. Поддерживает скачивание. Тип: info или error.
// @Tags         Логи
// @Security     BearerAuth
// @Produce      plain
// @Param        type     query string false "Тип логов: info (по умолчанию) или error"
// @Param        date     query string false "Дата в формате YYYY-MM-DD (по умолчанию — сегодня)"
// @Param        download query bool   false "Скачать файл (true) или отобразить в браузере"
// @Success      200 {string} string "Содержимое лог-файла"
// @Failure      400 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/admin/logs [get]
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	logType := query.Get("type")      // info | error
	date := query.Get("date")         // формат YYYY-MM-DD
	download := query.Get("download") // "true" => attachment

	if logType == "" {
		logType = "info"
	}
	validTypes := map[string]bool{"info": true, "error": true, "debug": true}
	if !validTypes[logType] {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid log type: must be 'info', 'error' or 'debug'")
		return
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	filePath := fmt.Sprintf("logs/%s.%s.log", logType, date)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ErrorJSON(w, http.StatusNotFound, "Log file not found for date: "+date)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		h.logger.Error("failed to open log file", zap.String("path", filePath), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Could not open log file")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain")
	if download == "true" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s.log", logType, date))
	}
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
