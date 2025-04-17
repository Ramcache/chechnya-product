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

type CartHandler struct {
	service services.CartServiceInterface
	logger  *zap.Logger
}

func NewCartHandler(service services.CartServiceInterface, logger *zap.Logger) *CartHandler {
	return &CartHandler{service: service, logger: logger}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (h *CartHandler) respondError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

// AddToCart
// @Summary Добавить товар в корзину
// @Description Добавляет выбранный товар в корзину пользователя
// @Tags Корзина
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param input body AddToCartRequest true "ID товара и количество"
// @Success 201 {string} string "Товар добавлен в корзину"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart [post]
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.service.AddToCart(userID, req.ProductID, req.Quantity); err != nil {
		h.logger.Error("add to cart failed", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("Item added to cart",
		zap.Int("user_id", userID),
		zap.Int("product_id", req.ProductID),
		zap.Int("quantity", req.Quantity),
	)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Added to cart"))
}

// GetCart
// @Summary Получить корзину
// @Description Возвращает список товаров в корзине текущего пользователя
// @Tags Корзина
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.CartItem
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	items, err := h.service.GetCart(userID)
	if err != nil {
		h.logger.Error("get cart failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "Failed to get cart")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UpdateItem
// @Summary Обновить количество товара в корзине
// @Description Изменяет количество указанного товара в корзине
// @Tags Корзина
// @Security BearerAuth
// @Accept json
// @Produce plain
// @Param product_id path int true "ID товара"
// @Param input body object{quantity=int} true "Новое количество"
// @Success 200 {string} string "Количество товара обновлено"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart/{product_id} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	productID, _ := strconv.Atoi(mux.Vars(r)["product_id"])
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := h.service.UpdateItem(userID, productID, req.Quantity); err != nil {
		h.logger.Warn("update item failed", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("Cart item updated",
		zap.Int("user_id", userID),
		zap.Int("product_id", productID),
		zap.Int("new_quantity", req.Quantity),
	)
	w.Write([]byte("Quantity updated"))
}

// DeleteItem
// @Summary Удалить товар из корзины
// @Description Удаляет указанный товар из корзины
// @Tags Корзина
// @Security BearerAuth
// @Produce plain
// @Param product_id path int true "ID товара"
// @Success 200 {string} string "Товар удалён из корзины"
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart/{product_id} [delete]
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	productID, _ := strconv.Atoi(mux.Vars(r)["product_id"])
	if err := h.service.DeleteItem(userID, productID); err != nil {
		h.logger.Error("delete item failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}
	h.logger.Info("Cart item deleted",
		zap.Int("user_id", userID),
		zap.Int("product_id", productID),
	)
	w.Write([]byte("Item deleted"))
}

// ClearCart
// @Summary Очистить корзину
// @Description Удаляет все товары из корзины пользователя
// @Tags Корзина
// @Security BearerAuth
// @Produce plain
// @Success 200 {string} string "Корзина очищена"
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart/clear [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if err := h.service.ClearCart(userID); err != nil {
		h.logger.Error("clear cart failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "Failed to clear cart")
		return
	}
	h.logger.Info("Cart cleared", zap.Int("user_id", userID))
	w.Write([]byte("Cart cleared"))
}

// Checkout
// @Summary Оформить заказ
// @Description Завершает оформление заказа на основе текущей корзины
// @Tags Корзина
// @Security BearerAuth
// @Produce plain
// @Success 200 {string} string "Заказ успешно оформлен"
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/cart/checkout [post]
func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if err := h.service.Checkout(userID); err != nil {
		h.logger.Error("checkout failed", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, "Checkout failed")
		return
	}
	h.logger.Info("Checkout completed",
		zap.Int("user_id", userID),
	)
	w.Write([]byte("Checkout successful"))
}
