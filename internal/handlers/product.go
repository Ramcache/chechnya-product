package handlers

import (
	"chechnya-product/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"chechnya-product/internal/services"
)

type ProductHandler struct {
	service *services.ProductService
	logger  *zap.Logger
}

func NewProductHandler(service *services.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{service: service, logger: logger}
}

// GetAll
// @Summary Получить список товаров
// @Description Получает список товаров с возможностью фильтрации
// @Tags Товар
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

	products, err := h.service.GetFiltered(search, category, minPrice, maxPrice, limit, offset, sort)
	if err != nil {
		h.logger.Error("failed to fetch products", zap.Error(err))
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	h.logger.Info("products fetched", zap.Int("count", len(products)))
	writeJSON(w, products)
}

// GetByID
// @Summary Получить товар по ID
// @Description Возвращает детали товара по его идентификатору
// @Tags Товар
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
		h.logger.Warn("invalid product ID", zap.String("id", idStr))
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Warn("product not found", zap.Int("id", id))
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	h.logger.Info("product fetched", zap.Int("id", id))
	writeJSON(w, product)
}

// Add
// @Summary Добавить товар (админ)
// @Description Создаёт новый товар (только для администратора)
// @Tags Товар
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
		h.logger.Warn("unauthorized access to add product")
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Warn("invalid product JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.AddProduct(&product); err != nil {
		h.logger.Error("failed to add product", zap.String("name", product.Name), zap.Error(err))
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	h.logger.Info("product added", zap.String("name", product.Name))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product added"))
}

// Update
// @Summary Обновить товар (админ)
// @Description Обновляет существующий товар по его ID (только для администратора)
// @Tags Товар
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
		h.logger.Warn("unauthorized access to update product")
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := parseIntParam(idStr)
	if err != nil {
		h.logger.Warn("invalid product ID for update", zap.String("id", idStr))
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Warn("invalid update JSON", zap.Error(err))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(id, &product); err != nil {
		h.logger.Error("failed to update product", zap.Int("id", id), zap.Error(err))
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	h.logger.Info("product updated", zap.Int("id", id))
	w.Write([]byte("Product updated"))
}

// Delete
// @Summary Удалить товар (админ)
// @Description Удаляет товар по его ID (только для администратора)
// @Tags Товар
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
		h.logger.Warn("unauthorized access to delete product")
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := parseIntParam(idStr)
	if err != nil {
		h.logger.Warn("invalid product ID for deletion", zap.String("id", idStr))
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		h.logger.Error("failed to delete product", zap.Int("id", id), zap.Error(err))
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	h.logger.Info("product deleted", zap.Int("id", id))
	w.Write([]byte("Product deleted"))
}

func parseIntParam(param string) (int, error) {
	return strconv.Atoi(param)
}
