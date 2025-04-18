package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const OwnerCookieName = "owner_id"
const OwnerHeaderName = "X-Owner-ID"

func GetOwnerID(w http.ResponseWriter, r *http.Request) string {
	// 1. Если пользователь авторизован
	userID := GetUserID(r)
	if userID != 0 {
		return "user_" + itoa(userID)
	}

	// 2. Попробовать из cookie
	if cookie, err := r.Cookie(OwnerCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// 3. Попробовать из заголовка
	if header := r.Header.Get(OwnerHeaderName); header != "" {
		return header
	}

	// 4. Сгенерировать guest ID
	guestID := "guest_" + uuid.New().String()

	// Установим cookie, чтобы сохранить
	http.SetCookie(w, &http.Cookie{
		Name:  OwnerCookieName,
		Value: guestID,
		Path:  "/",
	})

	return guestID
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
