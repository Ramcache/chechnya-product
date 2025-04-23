package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
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
// @Failure 400 {object} ErrorResponse
// @Router /api/cart [post]
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid AddToCart request", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.AddToCart(ownerID, req.ProductID, req.Quantity); err != nil {
		h.logger.Error("add to cart failed", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("item added to cart",
		zap.String("owner_id", ownerID),
		zap.Int("product_id", req.ProductID),
		zap.Int("quantity", req.Quantity),
	)
	JSONResponse(w, http.StatusCreated, "Added to cart", nil)

}

// GetCart
// @Summary Получить содержимое корзины
// @Description Возвращает список товаров в корзине для owner_id
// @Tags Корзина
// @Produce json
// @Success 200 {array} services.CartItemResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	items, err := h.service.GetCart(ownerID)
	if err != nil {
		h.logger.Error("get cart failed", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}

	if items == nil {
		items = []services.CartItemResponse{}
	}

	h.logger.Info("cart retrieved", zap.String("owner_id", ownerID), zap.Int("items_count", len(items)))
	JSONResponse(w, http.StatusOK, "Cart retrieved", items)
}

// UpdateItem
// @Summary Обновить количество товара в корзине
// @Description Обновляет количество указанного товара для owner_id
// @Tags Корзина
// @Accept json
// @Produce plain
// @Param product_id path int true "ID товара"
// @Param input body object{quantity=int} true "Новое количество"
// @Success 200 {string} string "Quantity updated"
// @Failure 400 {object} ErrorResponse
// @Router /api/cart/{product_id} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	productID, _ := strconv.Atoi(mux.Vars(r)["product_id"])
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid UpdateItem request", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := h.service.UpdateItem(ownerID, productID, req.Quantity); err != nil {
		h.logger.Warn("update item failed", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("cart item updated",
		zap.String("owner_id", ownerID),
		zap.Int("product_id", productID),
		zap.Int("new_quantity", req.Quantity),
	)
	JSONResponse(w, http.StatusOK, "Quantity updated", nil)
}

// DeleteItem
// @Summary Удалить товар из корзины
// @Description Удаляет товар по ID из корзины owner_id
// @Tags Корзина
// @Produce plain
// @Param product_id path int true "ID товара"
// @Success 200 {string} string "Item deleted"
// @Failure 500 {object} ErrorResponse
// @Router /api/cart/{product_id} [delete]
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	productID, _ := strconv.Atoi(mux.Vars(r)["product_id"])
	if err := h.service.DeleteItem(ownerID, productID); err != nil {
		h.logger.Error("delete item failed", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}
	h.logger.Info("cart item deleted",
		zap.String("owner_id", ownerID),
		zap.Int("product_id", productID),
	)

	JSONResponse(w, http.StatusOK, "Item deleted", nil)

}

// ClearCart
// @Summary Очистить корзину
// @Description Удаляет все товары из корзины owner_id
// @Tags Корзина
// @Produce plain
// @Success 200 {string} string "Cart cleared"
// @Failure 500 {object} ErrorResponse
// @Router /api/cart/clear [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	if err := h.service.ClearCart(ownerID); err != nil {
		h.logger.Error("clear cart failed", zap.Error(err), zap.String("owner_id", ownerID))
		ErrorJSON(w, http.StatusInternalServerError, "Failed to clear cart")
		return
	}
	h.logger.Info("cart cleared", zap.String("owner_id", ownerID))

	JSONResponse(w, http.StatusOK, "Cart cleared", nil)
}
