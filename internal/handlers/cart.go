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
