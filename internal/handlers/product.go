package handlers

import (
	"chechnya-product/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"chechnya-product/internal/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetAll godoc
// @Summary      Получить товары
// @Description  Возвращает товары с фильтрацией по поиску и категории
// @Tags         products
// @Produce      json
// @Param search   query string false "Поиск"
// @Param category query string false "Категория"
// @Param limit    query int    false "Сколько товаров вернуть"
// @Param offset   query int    false "С какого начать (смещение)"
// @Param sort     query string false "Сортировка: price_asc, price_desc, name_asc, name_desc, stock_asc, stock_desc"
// @Success      200      {array}   models.Product
// @Failure      500      {string}  string "Ошибка получения товаров"
// @Router       /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	category := r.URL.Query().Get("category")
	sort := r.URL.Query().Get("sort")

	// Разбираем limit
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Разбираем offset
	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	products, err := h.service.GetFiltered(search, category, limit, offset, sort)
	if err != nil {
		http.Error(w, "Ошибка получения товаров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Add godoc
// @Summary      Добавить новый товар
// @Description  Создаёт новый товар (доступно только администратору)
// @Tags         admin-products
// @Security     BearerAuth
// @Accept       json
// @Produce      plain
// @Param        product  body      models.Product  true  "Данные товара"
// @Success      201      {string}  string "Товар добавлен"
// @Failure      400      {string}  string "Невалидный JSON"
// @Failure      403      {string}  string "Нет доступа"
// @Failure      500      {string}  string "Ошибка добавления товара"
// @Router       /admin/products [post]
func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role")
	if role != "admin" {
		http.Error(w, "Доступ только для админов", http.StatusForbidden)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err := h.service.AddProduct(&product)
	if err != nil {
		http.Error(w, "Ошибка добавления товара", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Товар добавлен"))
}

// Delete godoc
// @Summary      Удалить товар
// @Description  Удаляет товар по ID (доступно только администратору)
// @Tags         admin-products
// @Security     BearerAuth
// @Produce      plain
// @Param        id   path      int  true  "ID товара"
// @Success      200  {string}  string "Товар удалён"
// @Failure      400  {string}  string "Некорректный ID"
// @Failure      403  {string}  string "Нет доступа"
// @Failure      500  {string}  string "Ошибка удаления товара"
// @Router       /admin/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role")
	if role != "admin" {
		http.Error(w, "Доступ только для админов", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, "Ошибка удаления товара", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Товар удалён"))
}

// Update godoc
// @Summary      Обновить товар
// @Description  Обновляет данные товара по ID (доступно только администратору)
// @Tags         admin-products
// @Security     BearerAuth
// @Accept       json
// @Produce      plain
// @Param        id       path      int            true  "ID товара"
// @Param        product  body      models.Product true  "Новые данные товара"
// @Success      200      {string}  string "Товар обновлён"
// @Failure      400      {string}  string "Некорректный ID или JSON"
// @Failure      500      {string}  string "Ошибка обновления товара"
// @Router       /admin/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, "Ошибка обновления товара", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Товар обновлён"))
}

// GetByID godoc
// @Summary      Получить товар по ID
// @Description  Возвращает один товар по его ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "ID товара"
// @Success      200  {object}  models.Product
// @Failure      404  {string}  string "Товар не найден"
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	product, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// GetCategories godoc
// @Summary      Получить список категорий
// @Description  Возвращает список всех уникальных категорий товаров
// @Tags         products
// @Produce      json
// @Success      200 {array} string
// @Failure      500 {string} string "Ошибка получения категорий"
// @Router       /categories [get]
func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		http.Error(w, "Ошибка получения категорий", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
