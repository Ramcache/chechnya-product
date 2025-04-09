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

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Регистрирует нового пользователя по логину и паролю
// @Tags         users
// @Accept       json
// @Produce      plain
// @Param        input  body      RegisterRequest  true  "Данные регистрации"
// @Success      201    {string}  string "Пользователь зарегистрирован"
// @Failure      400    {string}  string "Невалидный JSON или пользователь уже существует"
// @Router       /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err := h.service.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь зарегистрирован"))
}

// Login godoc
// @Summary      Авторизация
// @Description  Возвращает JWT токен при успешном входе
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body      LoginRequest  true  "Данные входа"
// @Success      200    {object}  map[string]string "JWT токен"
// @Failure      400    {string}  string "Невалидный JSON"
// @Failure      401    {string}  string "Неверный логин или пароль"
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cfg := config.LoadConfig()
	token, err := middleware.GenerateJWT(user.ID, user.Role, cfg.JWTSecret)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Me godoc
// @Summary      Профиль текущего пользователя
// @Description  Возвращает ID, имя пользователя и роль из JWT
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {string}  string "Пользователь не найден"
// @Router       /me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	user, err := h.service.GetByID(userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
