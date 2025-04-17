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

// GetAll
// @Summary Получить список товаров
// @Description Получает список товаров с возможностью фильтрации
// @Produce json
// @Param search query string false "Поиск по названию"
// @Param category query string false "Фильтр по категории"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param sort query string false "Сортировка"
// @Param limit query int false "Ограничение количества результатов"
// @Param offset query int false "Сдвиг для пагинации"
// @Success 200 {array} models.Product
// @Failure 500 {object} ErrorResponse
// @Router /api/products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	search := query.Get("search")
	category := query.Get("category")
	sort := query.Get("sort")

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	minPrice, _ := strconv.ParseFloat(query.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(query.Get("max_price"), 64)

	products, err := h.service.GetFiltered(
		search, category, minPrice, maxPrice,
		limit, offset, sort,
	)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	writeJSON(w, products)
}

// GetByID
// @Summary Получить товар по ID
// @Description Возвращает детали товара по его идентификатору
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {object} models.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := parseIntParam(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	writeJSON(w, product)
}

// Add
// @Summary Добавить товар (админ)
// @Description Создаёт новый товар (только для администратора)
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param input body models.Product true "Данные товара"
// @Success 201 {string} string "Product added"
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/products [post]
func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.AddProduct(&product); err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product added"))
}

// Update
// @Summary Обновить товар (админ)
// @Description Обновляет существующий товар по его ID (только для администратора)
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param id path int true "ID товара"
// @Param input body models.Product true "Новые данные товара"
// @Success 200 {string} string "Product updated"
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := parseIntParam(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(id, &product); err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Product updated"))
}

// Delete
// @Summary Удалить товар (админ)
// @Description Удаляет товар по его ID (только для администратора)
// @Security BearerAuth
// @Produce plain
// @Param id path int true "ID товара"
// @Success 200 {string} string "Product deleted"
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/admin/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := parseIntParam(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Product deleted"))
}

// GetCategories
// @Summary Получить категории товаров
// @Description Возвращает список категорий товаров
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} ErrorResponse
// @Router /api/categories [get]
func (h *ProductHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	writeJSON(w, categories)
}

func parseIntParam(param string) (int, error) {
	return strconv.Atoi(param)
}
