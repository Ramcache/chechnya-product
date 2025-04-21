package utils

import (
	"chechnya-product/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManagerInterface interface {
	Generate(userID int, role models.UserRole) (string, error)
	Verify(tokenStr string) (*UserClaims, error)
}

type JWTManager struct {
	secretKey string
	duration  time.Duration
}

// Claims, которые мы кладем в токен
type UserClaims struct {
	UserID int             `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// Конструктор JWT-менеджера
func NewJWTManager(secret string, duration time.Duration) *JWTManager {
	return &JWTManager{secretKey: secret, duration: duration}
}

// Генерация JWT токена
func (m *JWTManager) Generate(userID int, role models.UserRole) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Валидация и извлечение данных из токена
func (m *JWTManager) Verify(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
