package handlers

import (
	"chechnya-product/config"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
	logger  *zap.Logger
}

func NewUserHandler(service *services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
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
// @Tags Пользователь
// @Accept json
// @Produce plain
// @Param input body RegisterRequest true "Имя пользователя и пароль"
// @Success 201 {string} string "Пользователь успешно зарегистрирован"
// @Failure 400 {object} ErrorResponse
// @Router /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid register JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(req.Username, req.Password); err != nil {
		h.logger.Warn("registration failed", zap.String("username", req.Username), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("user registered", zap.String("username", req.Username))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User successfully registered"))
}

// Login
// @Summary Авторизация пользователя
// @Description Выполняет вход пользователя и возвращает JWT-токен
// @Tags Пользователь
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
		h.logger.Warn("invalid login JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		h.logger.Warn("login failed", zap.String("username", req.Username), zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cfg, _ := config.LoadConfig()
	token, err := middleware.GenerateJWT(user.ID, user.Role, cfg.JWTSecret)
	if err != nil {
		h.logger.Error("token generation failed", zap.Int("user_id", user.ID), zap.Error(err))
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	h.logger.Info("user logged in", zap.String("username", user.Username), zap.String("role", user.Role))
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Me
// @Summary Получить информацию о пользователе
// @Description Возвращает профиль текущего пользователя
// @Tags Пользователь
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
		h.logger.Warn("user not found", zap.Int("user_id", userID))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	h.logger.Info("user profile requested", zap.Int("user_id", user.ID), zap.String("username", user.Username))
	writeJSON(w, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
