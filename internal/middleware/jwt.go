package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(userID int, role string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenStr string, secret string) (int, string, error) {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		role := claims["role"].(string)
		return userID, role, nil
	}

	return 0, "", errors.New("invalid token")
}
