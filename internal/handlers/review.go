package handlers

import (
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ReviewHandlerInterface interface {
	AddReview(w http.ResponseWriter, r *http.Request)
	GetReviews(w http.ResponseWriter, r *http.Request)
	UpdateReview(w http.ResponseWriter, r *http.Request)
	DeleteReview(w http.ResponseWriter, r *http.Request)
}

type ReviewHandler struct {
	service services.ReviewServiceInterface
}

func NewReviewHandler(service services.ReviewServiceInterface) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// AddReview добавляет новый отзыв к товару
// @Summary Добавить отзыв
// @Description Отзыв может оставить как авторизованный, так и гость. Повторный отзыв от одного владельца невозможен.
// @Tags Отзывы
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Param review body models.ReviewRequest true "Оценка и комментарий"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/products/{id}/reviews [post]
func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)

	productID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var body struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if err := h.service.AddReview(ownerID, productID, body.Rating, body.Comment); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Review added", nil)
}

// GetReviews возвращает список отзывов по товару
// @Summary Получить отзывы товара
// @Tags Отзывы
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {array} models.Review
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id}/reviews [get]
func (h *ReviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(mux.Vars(r)["id"])
	reviews, err := h.service.GetReviewsByProductID(productID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Reviews fetched", reviews)
}

// UpdateReview обновляет отзыв по owner_id и product_id
// @Summary Обновить отзыв
// @Description Может обновить только тот, кто оставил (по owner_id)
// @Tags Отзывы
// @Accept json
// @Produce json
// @Param id path int true "ID товара"
// @Param review body models.ReviewRequest true "Обновлённая оценка и комментарий"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id}/reviews [put]
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	productID, _ := strconv.Atoi(mux.Vars(r)["id"])

	var body struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	err := h.service.UpdateReview(ownerID, productID, body.Rating, body.Comment)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to update review")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Review updated", nil)
}

// DeleteReview удаляет отзыв по owner_id и product_id
// @Summary Удалить отзыв
// @Description Может удалить только тот, кто оставил (по owner_id)
// @Tags Отзывы
// @Produce json
// @Param id path int true "ID товара"
// @Success 200 {object} utils.SuccessResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id}/reviews [delete]
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	productID, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := h.service.DeleteReview(ownerID, productID); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to delete review")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Review deleted", nil)
}
