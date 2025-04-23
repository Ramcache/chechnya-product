package middleware

import (
	"chechnya-product/internal/handlers"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					ownerID := GetOwnerID(w, r)

					// Логируем с полной трассировкой
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
						zap.String("owner_id", ownerID),
						zap.String("user_agent", r.UserAgent()),
						zap.String("stack", string(debug.Stack())),
					)

					handlers.ErrorJSON(w, http.StatusInternalServerError, "Internal Server Error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
