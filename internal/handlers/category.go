package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chechnya-product/internal/services"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAll
// @Summary Получить список категорий
// @Description Возвращает все доступные категории товаров
// @Tags Категории
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} ErrorResponse
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

// Create
// @Summary Создать новую категорию
// @Description Добавляет новую категорию (только для администратора)
// @Tags Категории
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param input body object{name=string} true "Название категории"
// @Success 201 {string} string "Category created"
// @Failure 400 {string} string "Invalid body or duplicate name"
// @Router /api/admin/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := h.service.Create(body.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Category created"))
}

// Update
// @Summary Обновить категорию
// @Description Изменяет название категории (только для администратора)
// @Tags Категории
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param id path int true "ID категории"
// @Param input body object{name=string} true "Новое имя категории"
// @Success 200 {string} string "Category updated"
// @Failure 400 {string} string "Invalid input or update failed"
// @Router /api/admin/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := h.service.Update(id, body.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Category updated"))
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Category deleted"))
}
