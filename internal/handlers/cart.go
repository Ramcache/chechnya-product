package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
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

// AddToCart godoc
// @Summary      Добавить товар в корзину
// @Description  Добавляет определённое количество товара в корзину пользователя
// @Tags         cart
// @Accept       json
// @Produce      plain
// @Param        input  body      AddToCartRequest  true  "Данные для добавления в корзину"
// @Success      201    {string}  string "Добавлено в корзину"
// @Failure      400    {string}  string "Невалидный запрос"
// @Failure      500    {string}  string "Ошибка добавления в корзину"
// @Router       /cart [post]
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

// GetCart godoc
// @Summary      Получить корзину пользователя
// @Description  Возвращает список товаров в корзине по user_id
// @Tags         cart
// @Produce      json
// @Param        user_id  query     int  true  "ID пользователя"
// @Success 200 {array} object
// @Failure      400      {string}  string "Неверный user_id"
// @Failure      500      {string}  string "Ошибка получения корзины"
// @Router       /cart [get]
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

// UpdateItem godoc
// @Summary      Обновить количество товара в корзине
// @Description  Обновляет количество определённого товара в корзине текущего пользователя
// @Tags         cart
// @Accept       json
// @Produce      plain
// @Param        product_id  path      int             true  "ID товара"
// @Param        input       body      map[string]int  true  "Новое количество (quantity)"
// @Success      200         {string}  string "Количество обновлено"
// @Failure      400         {string}  string "Невалидный JSON или ошибка"
// @Router       /cart/{product_id} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	productID, _ := strconv.Atoi(vars["product_id"])

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateItem(userID, productID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Количество обновлено"))
}

// DeleteItem godoc
// @Summary      Удалить товар из корзины
// @Description  Удаляет указанный товар из корзины текущего пользователя
// @Tags         cart
// @Produce      plain
// @Param        product_id  path      int  true  "ID товара"
// @Success      200         {string}  string "Товар удалён из корзины"
// @Failure      500         {string}  string "Ошибка удаления"
// @Router       /cart/{product_id} [delete]
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	productID, _ := strconv.Atoi(mux.Vars(r)["product_id"])

	err := h.service.DeleteItem(userID, productID)
	if err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Товар удалён из корзины"))
}
