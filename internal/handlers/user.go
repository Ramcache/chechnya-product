package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type UserHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Me(w http.ResponseWriter, r *http.Request)
	CreateUserByPhone(w http.ResponseWriter, r *http.Request)
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
// @Tags         Профиль
// @Accept       json
// @Param        register body RegisterRequest true "User registration data"
// @Produce json
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router       /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid register JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	oldOwnerID := middleware.GetOwnerID(w, r)

	user, err := h.service.Register(services.RegisterRequest{
		Phone:    req.Phone,
		Password: req.Password,
		Username: req.Username,
		Email:    req.Email,
		OwnerID:  oldOwnerID,
	})
	if err != nil {
		h.logger.Warn("registration failed", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if cartErr := h.service.TransferCart(oldOwnerID, user.OwnerID); cartErr != nil {
		h.logger.Warn("cart transfer failed", zap.String("from", oldOwnerID), zap.String("to", user.OwnerID), zap.Error(cartErr))
	}

	middleware.SetOwnerID(w, user.OwnerID)
	h.logger.Info("cart transferred",
		zap.String("from", oldOwnerID),
		zap.String("to", user.OwnerID),
	)

	h.logger.Info("user registered", zap.String("phone", user.Phone), zap.String("owner_id", user.OwnerID))
	utils.JSONResponse(w, http.StatusCreated, "Registration successful", nil)
}

// Login authenticates a user and returns JWT token
// @Summary      User login
// @Description  Logs in a user and returns JWT token on success
// @Tags         Профиль
// @Accept       json
// @Produce      json
// @Param        login body LoginRequest true "Phone and password"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router       /api/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid login JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	oldOwnerID := middleware.GetOwnerID(w, r) // 👈 получаем текущий owner_id (гость)

	user, token, err := h.service.LoginWithUser(services.LoginRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		h.logger.Warn("login failed", zap.String("identifier", req.Identifier), zap.Error(err))
		if err.Error() == "phone not verified" {
			utils.ErrorJSON(w, http.StatusForbidden, "Please verify your phone number first")
		} else {
			utils.ErrorJSON(w, http.StatusUnauthorized, "Invalid credentials")
		}
		return
	}

	if cartErr := h.service.TransferCart(oldOwnerID, user.OwnerID); cartErr != nil {
		h.logger.Warn("cart transfer failed", zap.String("from", oldOwnerID), zap.String("to", user.OwnerID), zap.Error(cartErr))
	}

	middleware.SetOwnerID(w, user.OwnerID)

	h.logger.Info("user logged in",
		zap.String("identifier", req.Identifier),
		zap.String("owner_id", user.OwnerID),
	)

	utils.JSONResponse(w, http.StatusOK, "Login successful", map[string]string{"token": token})
}

// Me returns current user profile
// @Summary      Get current user
// @Description  Returns profile info of authenticated user
// @Tags         Профиль
// @Security     BearerAuth
// @Produce      json
// @Success 200 {object} UserProfileResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      404 {string} string "User not found"
// @Router       /api/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.service.GetByID(claims.UserID)
	if err != nil || user == nil {
		h.logger.Warn("user not found", zap.Int("user_id", claims.UserID))
		utils.ErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}

	h.logger.Info("user profile requested", zap.Int("user_id", user.ID), zap.String("phone", user.Phone))
	utils.JSONResponse(w, http.StatusOK, "User profile", map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"isVerified": user.IsVerified,
		"owner_id":   user.OwnerID,
	})
}

type CreateByPhoneRequest struct {
	Phone string `json:"phone"`
}

// CreateUserByPhone создает пользователя по номеру телефона
// @Summary Создать пользователя по номеру
// @Tags Админ
// @Accept json
// @Produce json
// @Param input body CreateByPhoneRequest true "Номер телефона"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/admin/users [post]
func (h *UserHandler) CreateUserByPhone(w http.ResponseWriter, r *http.Request) {
	var req CreateByPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, password, err := h.service.CreateByPhone(req.Phone)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "User created", map[string]interface{}{
		"phone":    user.Phone,
		"password": password,
	})
}
