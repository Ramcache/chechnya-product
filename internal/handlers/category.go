package handlers

import (
	"chechnya-product/internal/utils"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"chechnya-product/internal/services"
	"github.com/gorilla/mux"
)

type CategoryHandlerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CategoryHandler struct {
	service services.CategoryServiceInterface
	logger  *zap.Logger
}

func NewCategoryHandler(service services.CategoryServiceInterface, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{service: service, logger: logger}
}

// GetAll
// @Summary Получить список категорий
// @Description Возвращает все доступные категории товаров
// @Tags Категории
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("failed to fetch categories", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}
	h.logger.Info("categories fetched", zap.Int("count", len(categories)))
	utils.JSONResponse(w, http.StatusOK, "Categories fetched", categories)
}

// Create
// @Summary Создать новую категорию
// @Description Добавляет новую категорию (только для администратора)
// @Tags Категории
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param input body utils.CategoryRequest true "Название категории"
// @Success 201 {string} string "Category created"
// @Failure 400 {string} string "Invalid body or duplicate name"
// @Router /api/admin/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		h.logger.Warn("invalid category creation request", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if err := h.service.Create(body.Name); err != nil {
		h.logger.Warn("failed to create category", zap.String("name", body.Name), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("category created", zap.String("name", body.Name))
	utils.JSONResponse(w, http.StatusCreated, "Category created", nil)
}

// Update
// @Summary Обновить категорию
// @Description Изменяет название категории (только для администратора)
// @Tags Категории
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор категории"
// @Param input body utils.CategoryRequest true "Новое имя категории"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/admin/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		h.logger.Warn("invalid update request", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if err := h.service.Update(id, body.Name); err != nil {
		h.logger.Warn("failed to update category", zap.Int("id", id), zap.String("name", body.Name), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("category updated", zap.Int("id", id), zap.String("name", body.Name))
	utils.JSONResponse(w, http.StatusOK, "Category updated", nil)
}

// Delete
// @Summary Удалить категорию
// @Description Удаляет категорию по ID (только для администратора)
// @Tags Категории
// @Security BearerAuth
// @Produce plain
// @Param id path int true "ID категории"
// @Success 200 {string} string "Category deleted"
// @Failure 400 {string} string "Удаление не удалось"
// @Router /api/admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(id); err != nil {
		h.logger.Error("failed to delete category", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("category deleted", zap.Int("id", id))
	utils.JSONResponse(w, http.StatusOK, "Category deleted", nil)
}
