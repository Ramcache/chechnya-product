package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type UserHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Me(w http.ResponseWriter, r *http.Request)
	CreateUserByPhone(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	service services.UserServiceInterface
	logger  *zap.Logger
}

func NewUserHandler(service services.UserServiceInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// Register — регистрация пользователя
// @Summary      Зарегистрировать нового пользователя
// @Description  Регистрирует нового пользователя по телефону, паролю, имени и e-mail
// @Tags         Профиль
// @Accept       json
// @Param        register body RegisterRequest true "Данные для регистрации"
// @Produce      json
// @Success      201 {object} utils.SuccessResponse
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Некорректный JSON при регистрации", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	oldOwnerID := middleware.GetOwnerID(w, r)

	if err := utils.ValidatePhone(req.Phone); err != nil {
		h.logger.Warn("Некорректный формат телефона", zap.String("phone", req.Phone))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(services.RegisterRequest{
		Phone:    req.Phone,
		Password: req.Password,
		Username: req.Username,
		Email:    req.Email,
		OwnerID:  oldOwnerID,
	})
	if err != nil {
		h.logger.Warn("Ошибка регистрации", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// переносим корзину если есть
	if cartErr := h.service.TransferCart(oldOwnerID, user.OwnerID); cartErr != nil {
		h.logger.Warn("Ошибка переноса корзины", zap.String("от", oldOwnerID), zap.String("к", user.OwnerID), zap.Error(cartErr))
	}

	middleware.SetOwnerID(w, user.OwnerID)
	h.logger.Info("Корзина перенесена",
		zap.String("от", oldOwnerID),
		zap.String("к", user.OwnerID),
	)

	h.logger.Info("Пользователь зарегистрирован", zap.String("phone", user.Phone), zap.String("owner_id", user.OwnerID))
	utils.JSONResponse(w, http.StatusCreated, "Регистрация успешна", nil)
}

// Login — аутентификация пользователя и выдача JWT
// @Summary      Вход пользователя
// @Description  Вход по телефону/почте и паролю. Возвращает JWT токен при успехе.
// @Tags         Профиль
// @Accept       json
// @Produce      json
// @Param        login body LoginRequest true "Телефон/почта и пароль"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Router       /api/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Некорректный JSON при входе", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	if err := utils.ValidateIdentifier(req.Identifier); err != nil {
		h.logger.Warn("Некорректный идентификатор", zap.String("identifier", req.Identifier), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	oldOwnerID := middleware.GetOwnerID(w, r)

	user, token, err := h.service.LoginWithUser(services.LoginRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		h.logger.Warn("Ошибка входа", zap.String("identifier", req.Identifier), zap.Error(err))
		if err.Error() == "phone not verified" {
			utils.ErrorJSON(w, http.StatusForbidden, "Сначала подтвердите номер телефона")
		} else {
			utils.ErrorJSON(w, http.StatusUnauthorized, "Неверные данные для входа")
		}
		return
	}

	if cartErr := h.service.TransferCart(oldOwnerID, user.OwnerID); cartErr != nil {
		h.logger.Warn("Ошибка переноса корзины", zap.String("от", oldOwnerID), zap.String("к", user.OwnerID), zap.Error(cartErr))
	}

	middleware.SetOwnerID(w, user.OwnerID)

	h.logger.Info("Пользователь вошёл",
		zap.String("identifier", req.Identifier),
		zap.String("owner_id", user.OwnerID),
	)

	utils.JSONResponse(w, http.StatusOK, "Вход выполнен успешно", map[string]string{
		"token":    token,
		"username": user.Username,
	})
}

// Me — получить профиль текущего пользователя
// @Summary      Получить профиль пользователя
// @Description  Возвращает данные профиля для авторизованного пользователя
// @Tags         Профиль
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} UserProfileResponse
// @Failure      401 {string} string "Не авторизован"
// @Failure      404 {string} string "Пользователь не найден"
// @Router       /api/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		utils.ErrorJSON(w, http.StatusUnauthorized, "Не авторизован")
		return
	}

	user, err := h.service.GetByID(claims.UserID)
	if err != nil || user == nil {
		h.logger.Warn("Пользователь не найден", zap.Int("user_id", claims.UserID))
		utils.ErrorJSON(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	h.logger.Info("Запрошен профиль пользователя", zap.Int("user_id", user.ID), zap.String("phone", user.Phone))
	utils.JSONResponse(w, http.StatusOK, "Профиль пользователя", map[string]interface{}{
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

// CreateUserByPhone — создать пользователя по номеру телефона
// @Summary Создать пользователя по номеру телефона
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
		utils.ErrorJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	user, password, err := h.service.CreateByPhone(req.Phone)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Пользователь создан", map[string]interface{}{
		"phone":    user.Phone,
		"password": password,
	})
}

// GetAllUsers возвращает список всех зарегистрированных пользователей
// @Summary Получить всех пользователей
// @Tags Профиль
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/users/all [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		h.logger.Error("не удалось получить пользователей", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка получения пользователей")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Пользователи получены", users)
}

// GetUserByID возвращает пользователя по ID
// @Summary Получить пользователя по ID
// @Tags Пользователи
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} utils.SuccessResponse{data=models.User}
// @Failure 400 {object} utils.ErrorResponse "Некорректный ID"
// @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// @Router /api/admin/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		h.logger.Warn("пользователь не найден", zap.Error(err))
		utils.ErrorJSON(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Пользователь найден", user)
}
