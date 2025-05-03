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
// @Summary      Получить лог-файл
// @Description  Возвращает содержимое лог-файла (info.log или error.log). Можно скачать файл или просмотреть в браузере.
// @Tags         Логи
// @Security     BearerAuth
// @Produce      plain
// @Param        type     query string false "Тип логов: info (по умолчанию) или error"
// @Param        download query bool   false "Скачать файл (true) или отобразить в браузере"
// @Success      200 {string} string "Содержимое лог-файла"
// @Failure      400 {object} utils.ErrorResponse "Неверный тип файла"
// @Failure      500 {object} utils.ErrorResponse "Ошибка при открытии файла"
// @Router       /api/admin/logs [get]
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	logType := query.Get("type")      // info | error
	download := query.Get("download") // "true" => attachment

	var filePath string
	switch logType {
	case "error":
		filePath = "logs/error.log"
	case "info", "":
		filePath = "logs/info.log"
	default:
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid log type")
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
		w.Header().Set("Content-Disposition", "attachment; filename="+logType+".log")
	}
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
