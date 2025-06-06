package handlers

import (
	"chechnya-product/internal/cache"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/models"
	"chechnya-product/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"chechnya-product/internal/services"
)

type ProductHandlerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	AddBulk(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	UploadImage(w http.ResponseWriter, r *http.Request)
}

type ProductHandler struct {
	service services.ProductServiceInterface
	logger  *zap.Logger
	cache   *cache.RedisCache
}

func NewProductHandler(service services.ProductServiceInterface, logger *zap.Logger, cache *cache.RedisCache) *ProductHandler {
	return &ProductHandler{service: service, logger: logger, cache: cache}
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
	availabilityStr := query.Get("availability")
	var availability *bool

	if availabilityStr == "true" {
		val := true
		availability = &val
	} else if availabilityStr == "false" {
		val := false
		availability = &val
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	minPrice, _ := strconv.ParseFloat(query.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(query.Get("max_price"), 64)

	cacheKey := fmt.Sprintf("products:search=%s:cat=%s:min=%.2f:max=%.2f:limit=%d:offset=%d:sort=%s:avail=%s",
		search, category, minPrice, maxPrice, limit, offset, sort, availabilityStr)

	ctx := r.Context()
	var cached []models.ProductCache

	err := h.cache.GetOrSet(ctx, cacheKey, &cached, func() (any, error) {
		products, err := h.service.GetFilteredRaw(search, category, minPrice, maxPrice, limit, offset, sort, availability)
		if err != nil {
			return nil, err
		}
		return models.ConvertProductsToCache(products), nil
	})
	if err != nil {
		h.logger.Error("cache fetch failed", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	products := models.ConvertCacheToProducts(cached)

	var result []models.ProductResponse
	for _, p := range products {
		var categoryName string
		if p.CategoryID.Valid {
			name, err := h.service.GetCategoryNameByID(int(p.CategoryID.Int64))
			if err == nil {
				categoryName = name
			}
		}
		response := utils.BuildProductResponse(&p, categoryName)
		result = append(result, response)
	}

	h.logger.Info("products fetched (cached or fresh)", zap.Int("count", len(result)))
	utils.JSONResponse(w, http.StatusOK, "Products fetched", result)
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

	ctx := r.Context()
	cacheKey := fmt.Sprintf("product:%d", id)

	var product models.ProductResponse

	err = h.cache.GetOrSet(ctx, cacheKey, &product, func() (any, error) {
		return h.service.GetByID(id)
	})

	h.logger.Info("product fetched (cached or fresh)", zap.Int("id", id))
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

	var input models.ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Warn("invalid product JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	product := mapProductInputToProduct(input)

	if err := h.service.AddProduct(&product); err != nil {
		h.logger.Error("failed to add product", zap.String("name", product.Name), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to add product")
		return
	}

	h.cache.ClearPrefix(r.Context(), "products:")

	response := utils.BuildProductResponse(&product, "")
	h.logger.Info("product added", zap.String("name", product.Name))
	utils.JSONResponse(w, http.StatusCreated, "Product added", response)
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
		h.logger.Warn("unauthorized access to update product")
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

	var input models.ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Warn("invalid update JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	product := mapProductInputToProduct(input)

	response, err := h.service.UpdateProduct(id, &product)
	if err != nil {
		h.logger.Error("failed to update product", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	h.cache.ClearPrefix(r.Context(), "products:")
	h.cache.Delete(r.Context(), fmt.Sprintf("product:%d", id))

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
	h.cache.ClearPrefix(r.Context(), "products:")
	h.cache.Delete(r.Context(), fmt.Sprintf("product:%d", id))

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

	var inputs []models.ProductInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil || len(inputs) == 0 {
		h.logger.Warn("invalid bulk product JSON", zap.Error(err))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON or empty array")
		return
	}

	var products []models.Product
	for _, input := range inputs {
		products = append(products, mapProductInputToProduct(input))
	}

	responses, err := h.service.AddProductsBulk(products)
	if err != nil {
		h.logger.Error("failed to bulk add products", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to add products")
		return
	}
	h.cache.ClearPrefix(r.Context(), "products:")

	h.logger.Info("bulk products added", zap.Int("count", len(responses)))
	utils.JSONResponse(w, http.StatusCreated, "Products added", responses)
}

// Patch
// @Summary Частичное обновление товара (админ)
// @Description Обновляет отдельные поля товара по его ID (только для администратора)
// @Tags Товар
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Param input body models.ProductPatchInput true "Поля для обновления товара"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/products/{id} [patch]
func (h *ProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil || claims.Role != "admin" {
		utils.ErrorJSON(w, http.StatusForbidden, "Access denied")
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := utils.ParseIntParam(idStr)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var input models.ProductPatchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Price != nil {
		if *input.Price <= 0 {
			utils.ErrorJSON(w, http.StatusBadRequest, "Price must be positive")
			return
		}
		updates["price"] = *input.Price
	}
	if input.Availability != nil {
		updates["availability"] = *input.Availability
	}
	if input.CategoryID != nil {
		updates["category_id"] = *input.CategoryID
	}
	if input.Url != nil {
		updates["url"] = *input.Url
	}

	if len(updates) == 0 {
		utils.ErrorJSON(w, http.StatusBadRequest, "No fields to update")
		return
	}

	if err := h.service.PatchProduct(id, updates); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	h.cache.ClearPrefix(r.Context(), "products:")
	h.cache.Delete(r.Context(), fmt.Sprintf("product:%d", id))

	utils.JSONResponse(w, http.StatusOK, "Product updated", nil)
}

func mapProductInputToProduct(input models.ProductInput) models.Product {
	availability := true
	if input.Availability != nil {
		availability = *input.Availability
	}

	product := models.Product{
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		Availability: availability,
		Url:          sql.NullString{String: input.Url, Valid: input.Url != ""},
	}

	if input.CategoryID != nil {
		product.CategoryID = sql.NullInt64{Int64: int64(*input.CategoryID), Valid: true}
	} else {
		product.CategoryID = sql.NullInt64{Valid: false}
	}

	return product
}

// UploadImage загружает изображение и возвращает ссылку
// @Summary Загрузить изображение
// @Tags Загрузка
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Файл изображения"
// @Success 200 {object} utils.SuccessResponse{data=map[string]string}
// @Failure 400 {object} utils.ErrorResponse
// @Router /admin/api/upload [post]
func (h *ProductHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Парсим multipart
	file, header, err := r.FormFile("image")
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Файл не получен")
		return
	}
	defer file.Close()

	// Генерируем уникальное имя
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	savePath := filepath.Join("uploads", filename)

	// Создаём файл
	out, err := os.Create(savePath)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Не удалось сохранить файл")
		return
	}
	defer out.Close()

	// Копируем содержимое
	_, err = io.Copy(out, file)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Ошибка при сохранении файла")
		return
	}

	// Ссылка на файл
	url := fmt.Sprintf("https://chechnya-product.ru/uploads/%s", filename)

	utils.JSONResponse(w, http.StatusOK, "Изображение загружено", map[string]string{
		"url": url,
	})
}
