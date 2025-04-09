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

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения товаров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role")
	if role != "admin" {
		http.Error(w, "Доступ только для админов", http.StatusForbidden)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err := h.service.AddProduct(&product)
	if err != nil {
		http.Error(w, "Ошибка добавления товара", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Товар добавлен"))
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role")
	if role != "admin" {
		http.Error(w, "Доступ только для админов", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, "Ошибка удаления товара", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Товар удалён"))
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, "Ошибка обновления товара", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Товар обновлён"))
}
