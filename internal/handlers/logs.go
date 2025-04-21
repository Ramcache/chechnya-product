package handlers

import (
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

// GetLogs — возвращает содержимое лог-файла
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(h.logPath)
	if err != nil {
		h.logger.Error("failed to open log file", zap.Error(err))
		http.Error(w, "Could not read logs", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain")
	io.Copy(w, file)
}
