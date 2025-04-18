package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserID   contextKey = "userID"
	ContextUserRole contextKey = "userRole"
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
			ctx = context.WithValue(ctx, ContextUserRole, role)
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

func GetUserRole(r *http.Request) string {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		return ""
	}
	return role
}

func GetUserIDOrZero(r *http.Request) int {
	id, ok := r.Context().Value("user_id").(int)
	if !ok {
		return 0
	}
	return id
}
