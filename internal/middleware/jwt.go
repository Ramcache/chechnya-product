package middleware

import (
	"chechnya-product/internal/utils"
	"context"
	"net/http"
	"strings"
)

// JWTMiddleware добавляет claims в context
func JWTMiddleware(jwt utils.JWTManagerInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := jwt.Verify(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OnlyAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetUserClaims(r)
			if claims == nil || claims.Role != "admin" {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDOrZero(r *http.Request) int {
	if claims := GetUserClaims(r); claims != nil {
		return claims.UserID
	}
	return 0
}
