package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			userID, role, err := ParseJWT(tokenStr, secret)
			if err != nil {
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			ctx = context.WithValue(ctx, ContextRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Получение userID из контекста
func GetUserID(r *http.Request) int {
	if val, ok := r.Context().Value(ContextUserID).(int); ok {
		return val
	}
	return 0
}
