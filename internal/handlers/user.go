package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type UserHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Me(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	service services.UserServiceInterface
	logger  *zap.Logger
}

func NewUserHandler(service services.UserServiceInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// Register handles user registration
// @Summary      Register new user
// @Description  Registers a new user with phone, password, username and email
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Param        register body RegisterRequest true "User registration data"
// @Success      201 {string} string "Registration successful"
// @Failure      400 {string} string "Invalid JSON or registration error"
// @Router       /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid register JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 👇 сохраняем старый guest ID до регистрации
	oldOwnerID := middleware.GetOwnerID(w, r)

	// 🧾 создаём пользователя
	user, err := h.service.Register(services.RegisterRequest{
		Phone:    req.Phone,
		Password: req.Password,
		Username: req.Username,
		Email:    req.Email,
		OwnerID:  oldOwnerID,
	})
	if err != nil {
		h.logger.Warn("registration failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ✅ переносим корзину на новый owner_id
	newOwnerID := user.OwnerID
	if cartErr := h.service.TransferCart(oldOwnerID, newOwnerID); cartErr != nil {
		h.logger.Warn("cart transfer failed", zap.String("from", oldOwnerID), zap.String("to", newOwnerID), zap.Error(cartErr))
	}

	// 🧹 удаляем старый owner_id из куки
	http.SetCookie(w, &http.Cookie{
		Name:   middleware.OwnerCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	h.logger.Info("user registered", zap.String("phone", user.Phone), zap.String("owner_id", newOwnerID))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registration successful"))
}

// Login authenticates a user and returns JWT token
// @Summary      User login
// @Description  Logs in a user and returns JWT token on success
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body LoginRequest true "Phone and password"
// @Success      200 {object} map[string]string "token"
// @Failure      400 {string} string "Invalid JSON"
// @Failure      401 {string} string "Invalid credentials"
// @Router       /api/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid login JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(services.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Warn("login failed", zap.String("phone", req.Phone), zap.Error(err))

		if err.Error() == "phone not verified" {
			http.Error(w, "Please verify your phone number first", http.StatusForbidden)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
		return
	}

	h.logger.Info("user logged in", zap.String("phone", req.Phone))
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Me returns current user profile
// @Summary      Get current user
// @Description  Returns profile info of authenticated user
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      401 {string} string "Unauthorized"
// @Failure      404 {string} string "User not found"
// @Router       /api/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.service.GetByID(claims.UserID)
	if err != nil || user == nil {
		h.logger.Warn("user not found", zap.Int("user_id", claims.UserID))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	h.logger.Info("user profile requested", zap.Int("user_id", user.ID), zap.String("phone", user.Phone))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"isVerified": user.IsVerified,
		"owner_id":   user.OwnerID,
	})
}
