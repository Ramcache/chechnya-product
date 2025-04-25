package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/models"
	"chechnya-product/internal/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"chechnya-product/internal/services"
)

type ProductHandlerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	AddBulk(w http.ResponseWriter, r *http.Request)
}

type ProductHandler struct {
	service services.ProductServiceInterface
	logger  *zap.Logger
}

func NewProductHandler(service services.ProductServiceInterface, logger *zap.Logger) *ProductHandler {
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
// @Success 200 {array} models.ProductResponse
// @Failure 500 {object} utils.ErrorResponse
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
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	h.logger.Info("products fetched", zap.Int("count", len(products)))
	utils.JSONResponse(w, http.StatusOK, "Products fetched", products)
}

// GetByID
// @Summary Получить товар по ID
// @Description Возвращает детали товара по его идентификатору
// @Tags Товар
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := utils.ParseIntParam(idStr)
	if err != nil {
		h.logger.Warn("invalid product ID", zap.String("id", idStr))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Warn("product not found", zap.Int("id", id))
		utils.ErrorJSON(w, http.StatusNotFound, "Product not found")
		return
	}

	h.logger.Info("product fetched", zap.Int("id", id))
	utils.JSONResponse(w, http.StatusOK, "Product fetched", product)
}

// Add
// @Summary Добавить товар (админ)
// @Description Создаёт новый товар (только для администратора)
// @Tags Товар
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Product true "Данные товара"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/products [post]
func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil || claims.Role != "admin" {
		h.logger.Warn("unauthorized access to add product")
		utils.ErrorJSON(w, http.StatusForbidden, "Access denied")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Warn("invalid product JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.service.AddProduct(&product); err != nil {
		h.logger.Error("failed to add product", zap.String("name", product.Name), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to add product")
		return
	}

	h.logger.Info("product added", zap.String("name", product.Name))
	utils.JSONResponse(w, http.StatusCreated, "Product added", nil)
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
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil || claims.Role != "admin" {
		h.logger.Warn("unauthorized access to add product")
		utils.ErrorJSON(w, http.StatusForbidden, "Access denied")
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := utils.ParseIntParam(idStr)
	if err != nil {
		h.logger.Warn("invalid product ID for update", zap.String("id", idStr))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.logger.Warn("invalid update JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	response, err := h.service.UpdateProduct(id, &product)
	if err != nil {
		h.logger.Error("failed to update product", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	h.logger.Info("product updated", zap.Int("id", id))
	utils.JSONResponse(w, http.StatusOK, "Product updated", response)
}

// Delete
// @Summary Удалить товар (админ)
// @Description Удаляет товар по его ID (только для администратора)
// @Tags Товар
// @Security BearerAuth
// @Produce plain
// @Param id path int true "ID товара"
// @Success 200 {string} string "Product deleted"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil || claims.Role != "admin" {
		h.logger.Warn("unauthorized access to add product")
		utils.ErrorJSON(w, http.StatusForbidden, "Access denied")
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := utils.ParseIntParam(idStr)
	if err != nil {
		h.logger.Warn("invalid product ID for deletion", zap.String("id", idStr))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		h.logger.Error("failed to delete product", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	h.logger.Info("product deleted", zap.Int("id", id))
	utils.JSONResponse(w, http.StatusOK, "Product deleted", nil)
}

// AddBulk
// @Summary Массовое добавление товаров (админ)
// @Description Добавляет несколько товаров сразу (только для администратора)
// @Tags Товар
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body []models.Product true "Массив товаров"
// @Success 201 {array} models.ProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/products/bulk [post]
func (h *ProductHandler) AddBulk(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil || claims.Role != "admin" {
		h.logger.Warn("unauthorized access to bulk add products")
		utils.ErrorJSON(w, http.StatusForbidden, "Access denied")
		return
	}

	var products []models.Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil || len(products) == 0 {
		h.logger.Warn("invalid bulk product JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON or empty array")
		return
	}

	responses, err := h.service.AddProductsBulk(products)
	if err != nil {
		h.logger.Error("failed to bulk add products", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to add products")
		return
	}

	h.logger.Info("bulk products added", zap.Int("count", len(responses)))
	utils.JSONResponse(w, http.StatusCreated, "Products added", responses)

}
