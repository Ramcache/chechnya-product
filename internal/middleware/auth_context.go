package middleware

import (
	"chechnya-product/internal/utils"
	"net/http"
)

// userClaimsKey используется для извлечения данных из context
// (тот же ключ используется в JWTMiddleware)
type contextKey string

const userClaimsKey contextKey = "user"

// GetUserClaims возвращает claims из контекста, если есть
func GetUserClaims(r *http.Request) *utils.UserClaims {
	if claims, ok := r.Context().Value(userClaimsKey).(*utils.UserClaims); ok {
		return claims
	}
	return nil
}

// GetUserID возвращает ID пользователя из JWT или 0, если не авторизован
func GetUserID(r *http.Request) int {
	claims := GetUserClaims(r)
	if claims != nil {
		return claims.UserID
	}
	return 0
}

// GetUserRole возвращает роль пользователя ("user", "admin"), если есть
func GetUserRole(r *http.Request) string {
	claims := GetUserClaims(r)
	if claims != nil {
		return string(claims.Role)
	}
	return "guest"
}

// IsAuthenticated проверяет, авторизован ли пользователь
func IsAuthenticated(r *http.Request) bool {
	return GetUserClaims(r) != nil
}

// IsAdmin проверяет, является ли пользователь администратором
func IsAdmin(r *http.Request) bool {
	return GetUserRole(r) == "admin"
}
