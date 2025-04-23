package middleware

import (
	"chechnya-product/internal/handlers"
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
				handlers.ErrorJSON(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := jwt.Verify(tokenStr)
			if err != nil {
				handlers.ErrorJSON(w, http.StatusUnauthorized, "Invalid token")
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
				handlers.ErrorJSON(w, http.StatusForbidden, "Access denied")
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
