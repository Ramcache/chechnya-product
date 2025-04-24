package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const OwnerCookieName = "owner_id"
const OwnerHeaderName = "X-Owner-ID"

// GetOwnerID определяет ID владельца корзины: user или guest.
// Всегда сохраняет owner_id в cookie, если он был получен.
func GetOwnerID(w http.ResponseWriter, r *http.Request) string {
	// 1. Если пользователь авторизован — user_x
	userID := GetUserID(r)
	if userID != 0 {
		ownerID := "user_" + itoa(userID)
		setOwnerCookie(w, ownerID)
		return ownerID
	}

	// 2. Из cookie
	if cookie, err := r.Cookie(OwnerCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// 3. Из заголовка
	if header := r.Header.Get(OwnerHeaderName); header != "" {
		setOwnerCookie(w, header) // Сохраним в cookie для следующих запросов
		return header
	}

	// 4. Новый guest ID
	guestID := "guest_" + uuid.New().String()
	setOwnerCookie(w, guestID)
	return guestID
}

// setOwnerCookie сохраняет owner_id в cookie
func setOwnerCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  OwnerCookieName,
		Value: value,
		Path:  "/",
	})
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
