package middleware

import (
	"net/http"
)

// Middleware: доступ только для admin
func OnlyAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextRole).(string)
			if !ok || role != "admin" {
				http.Error(w, "Доступ запрещён: только для администраторов", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
