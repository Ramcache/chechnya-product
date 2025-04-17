package handlers

import (
	"chechnya-product/config"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя
// @Accept json
// @Produce plain
// @Param input body RegisterRequest true "Имя пользователя и пароль"
// @Success 201 {string} string "Пользователь успешно зарегистрирован"
// @Failure 400 {object} ErrorResponse
// @Router /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(req.Username, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User successfully registered"))
}

// Login
// @Summary Авторизация пользователя
// @Description Выполняет вход пользователя и возвращает JWT-токен
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Имя пользователя и пароль"
// @Success 200 {object} object{token=string}
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cfg, _ := config.LoadConfig()

	token, err := middleware.GenerateJWT(user.ID, user.Role, cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Me
// @Summary Получить информацию о пользователе
// @Description Возвращает профиль текущего пользователя
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{id=int, username=string, role=string}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	user, err := h.service.GetByID(userID)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
