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
	CreateBulk(w http.ResponseWriter, r *http.Request)
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
// @Produce json
// @Param input body utils.CategoryRequest true "Название категории"
// @Success 201 {object} models.Category
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/admin/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body utils.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		h.logger.Warn("invalid category creation request", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	category, err := h.service.Create(body.Name, body.SortOrder)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSONResponse(w, http.StatusCreated, "Category created", category)

	h.logger.Info("category created", zap.String("name", category.Name), zap.Int("sortOrder", category.SortOrder))
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
		Name      *string `json:"name"`
		SortOrder *int    `json:"sortOrder"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Warn("invalid update request", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if body.Name == nil && body.SortOrder == nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Nothing to update")
		return
	}

	updatedCategory, err := h.service.PartialUpdate(id, body.Name, body.SortOrder)
	if err != nil {
		h.logger.Warn("failed to update category", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("category updated", zap.Int("id", id))
	utils.JSONResponse(w, http.StatusOK, "Category updated", updatedCategory)
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

// CreateBulk
// @Summary Массовое создание категорий
// @Description Добавляет несколько категорий сразу (только для администратора)
// @Tags Категории
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body []utils.CategoryRequest true "Список категорий"
// @Success 201 {string} string "Categories created"
// @Failure 400 {string} string "Invalid body"
// @Router /api/admin/categories/bulk [post]
func (h *CategoryHandler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var categories []utils.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&categories); err != nil || len(categories) == 0 {
		h.logger.Warn("invalid bulk create request", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body or empty array")
		return
	}

	created, err := h.service.CreateBulk(categories)
	if err != nil {
		h.logger.Error("bulk category creation failed", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to create categories")
		return
	}

	names := make([]string, 0, len(created))
	for _, c := range created {
		names = append(names, c.Name)
	}
	h.logger.Info("bulk categories created", zap.Int("count", len(created)), zap.Strings("names", names))
	utils.JSONResponse(w, http.StatusCreated, "Categories created", created)
}
