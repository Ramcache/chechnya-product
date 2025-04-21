package handlers

// Запрос на регистрацию
type RegisterRequest struct {
	Phone    string  `json:"phone"`
	Password string  `json:"password"`
	Username string  `json:"username"`
	Email    *string `json:"email,omitempty"`
}

// Запрос на вход
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
