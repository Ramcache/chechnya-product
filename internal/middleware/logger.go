package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Получаем owner_id (user_123 или guest_xxx)
			ownerID := GetOwnerID(w, r)

			next.ServeHTTP(w, r)

			// Логируем с owner_id
			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.String("owner_id", ownerID),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
