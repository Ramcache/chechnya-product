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
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserProfileResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	IsVerified bool   `json:"isVerified"`
	OwnerID    string `json:"owner_id"`
}
