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
// @Description  Возвращает лог за указанную дату. Поддерживает скачивание.
// @Tags         Логи
// @Security     BearerAuth
// @Produce      plain
// @Param        date     query string false "Дата в формате YYYY-MM-DD (по умолчанию — сегодня)"
// @Param        download query bool   false "Скачать файл (true) или отобразить в браузере"
// @Success      200 {string} string "Содержимое лог-файла"
// @Failure      400 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/admin/logs [get]
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	download := query.Get("download")
	date := query.Get("date") // формат YYYY-MM-DD

	if date == "" {
		date = utils.TodayDate() // utils.TodayDate() вернёт текущую дату в формате "2006-01-02"
	}

	// Формируем путь к файлу
	filePath := "logs/app." + date + ".log"

	// Проверка, существует ли файл
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
		w.Header().Set("Content-Disposition", "attachment; filename=log-"+date+".log")
	}

	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
