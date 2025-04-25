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

// POST /api/products/{id}/reviews
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

// GET /api/products/{id}/reviews
func (h *ReviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(mux.Vars(r)["id"])
	reviews, err := h.service.GetReviewsByProductID(productID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Reviews fetched", reviews)
}

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

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	ownerID := middleware.GetOwnerID(w, r)
	productID, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err := h.service.DeleteReview(ownerID, productID); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to delete review")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Review deleted", nil)
}
