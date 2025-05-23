package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/models"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CartHandlerInterface interface {
	AddToCart(w http.ResponseWriter, r *http.Request)
	GetCart(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	ClearCart(w http.ResponseWriter, r *http.Request)
	AddBulkToCart(w http.ResponseWriter, r *http.Request)
}

type CartHandler struct {
	service services.CartServiceInterface
	logger  *zap.Logger
}

func NewCartHandler(service services.CartServiceInterface, logger *zap.Logger) *CartHandler {
	return &CartHandler{service: service, logger: logger}
}

type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// AddToCart
// @Summary Добавить товар в корзину
// @Description Добавляет товар в корзину по owner_id (user или ip)
// @Tags Корзина
// @Accept json
// @Produce plain
// @Param input body AddToCartRequest true "ID товара и количество"
// @Success 201 {string} string "Added to cart"
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/cart [post]
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid AddToCart request", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AddToCart(ownerID, req.ProductID, req.Quantity); err != nil {
		h.logger.Error("add to cart failed", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	cartItems, err := h.service.GetCart(ownerID)
	if err != nil {
		h.logger.Error("get cart after add failed", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch updated cart")
		return
	}

	var addedItem *models.CartItemResponse
	for _, item := range cartItems {
		if item.ProductID == req.ProductID {
			addedItem = &item
			break
		}
	}

	if addedItem == nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Item added but not found in cart")
		return
	}

	h.logger.Info("item added to cart",
		zap.String("owner_id", ownerID),
		zap.Int("product_id", req.ProductID),
		zap.Int("quantity", req.Quantity),
	)

	utils.JSONResponse(w, http.StatusCreated, "Item added to cart", addedItem)
}

// GetCart
// @Summary Получить содержимое корзины
// @Description Возвращает список товаров в корзине для owner_id
// @Tags Корзина
// @Produce json
// @Success 200 {object} models.CartBulkResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	items, err := h.service.GetCart(ownerID)
	if err != nil {
		h.logger.Error("get cart failed", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	if items == nil {
		items = []models.CartItemResponse{}
	}

	h.logger.Info("cart retrieved", zap.String("owner_id", ownerID), zap.Int("items_count", len(items)))
	var total float64
	for _, item := range items {
		total += item.Total
	}

	response := models.CartBulkResponse{
		Items: items,
		Total: total,
	}

	utils.JSONResponse(w, http.StatusOK, "Cart retrieved", response)

}

// UpdateItem
// @Summary Обновить количество товара в корзине
// @Description Обновляет количество указанного товара для owner_id
// @Tags Корзина
// @Accept json
// @Produce plain
// @Param product_id path int true "ID товара"
// @Param input body utils.UpdateItemRequest true "Новое количество"
// @Success 200 {string} string "Quantity updated"
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/cart/{product_id} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		h.logger.Warn("invalid product_id", zap.String("raw", productIDStr))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid UpdateItem request", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := h.service.UpdateItem(ownerID, productID, req.Quantity); err != nil {
		h.logger.Warn("update item failed", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("cart item updated",
		zap.String("owner_id", ownerID),
		zap.Int("product_id", productID),
		zap.Int("new_quantity", req.Quantity),
	)
	utils.JSONResponse(w, http.StatusOK, "Quantity updated", nil)
}

// DeleteItem
// @Summary Удалить товар из корзины
// @Description Удаляет товар по ID из корзины owner_id
// @Tags Корзина
// @Produce plain
// @Param product_id path int true "Идентификатор товара, который нужно удалить или обновить"
// @Success 200 {string} string "Item deleted"
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/cart/{product_id} [delete]
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		h.logger.Warn("invalid product_id", zap.String("raw", productIDStr))
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.DeleteItem(ownerID, productID); err != nil {
		h.logger.Error("delete item failed", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	h.logger.Info("cart item deleted",
		zap.String("owner_id", ownerID),
		zap.Int("product_id", productID),
	)

	utils.JSONResponse(w, http.StatusOK, "Item deleted", nil)
}

// ClearCart
// @Summary Очистить корзину
// @Description Удаляет все товары из корзины owner_id
// @Tags Корзина
// @Produce plain
// @Success 200 {string} string "Cart cleared"
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/cart/clear [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	if err := h.service.ClearCart(ownerID); err != nil {
		h.logger.Error("clear cart failed", zap.Error(err), zap.String("owner_id", ownerID))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to clear cart")
		return
	}
	h.logger.Info("cart cleared", zap.String("owner_id", ownerID))

	utils.JSONResponse(w, http.StatusOK, "Cart cleared", nil)
}

// AddBulkToCart
// @Summary Добавить несколько товаров в корзину
// @Description Добавляет несколько товаров в корзину за один запрос (bulk).
// @Tags Корзина
// @Accept json
// @Produce json
// @Param input body []AddToCartRequest true "Список товаров для добавления"
// @Success 201 {object} utils.SuccessResponse{data=models.CartBulkResponse} "Товары добавлены, возвращены список и сумма"
// @Failure 400 {object} utils.ErrorResponse "Некорректные данные запроса"
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера при получении корзины"
// @Router /api/cart/bulk [post]
func (h *CartHandler) AddBulkToCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	var items []AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for _, item := range items {
		if item.Quantity <= 0 {
			continue
		}
		err := h.service.AddToCart(ownerID, item.ProductID, item.Quantity)
		if err != nil {
			h.logger.Warn("bulk add failed",
				zap.Int("product_id", item.ProductID),
				zap.Error(err),
			)
		}
	}

	cartItems, err := h.service.GetCart(ownerID)
	if err != nil {
		h.logger.Error("get cart after bulk failed", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to get updated cart")
		return
	}

	var total float64
	for _, item := range cartItems {
		total += item.Total
	}

	h.logger.Info("bulk items added", zap.String("owner_id", ownerID), zap.Int("count", len(cartItems)))

	utils.JSONResponse(w, http.StatusCreated, "Items added to cart", models.CartBulkResponse{
		Items: cartItems,
		Total: total,
	})

}
