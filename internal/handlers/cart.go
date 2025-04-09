package handlers

import (
	"chechnya-product/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

type AddToCartRequest struct {
	UserID    int `json:"user_id"` // временно напрямую
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	err := h.service.AddToCart(req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, "Ошибка добавления в корзину", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Добавлено в корзину"))
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный user_id", http.StatusBadRequest)
		return
	}

	items, err := h.service.GetCart(userID)
	if err != nil {
		http.Error(w, "Ошибка получения корзины", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
