package utils

import (
	"chechnya-product/internal/models"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

type CategoryRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if message == "" {
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	})
}
func ErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func ParseIntParam(param string) (int, error) {
	return strconv.Atoi(param)
}

func BuildProductResponse(p *models.Product, categoryName string) models.ProductResponse {
	var categoryID int
	if p.CategoryID.Valid {
		categoryID = int(p.CategoryID.Int64)
	} else {
		categoryID = 0
	}

	var url string
	if p.Url.Valid {
		url = p.Url.String
	} else {
		url = ""
	}

	return models.ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Price:        p.Price,
		Availability: p.Availability,
		CategoryID:   categoryID,
		CategoryName: categoryName,
		Url:          url,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#"
const idCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

func GeneratePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	pass := make([]byte, length)
	for i := range pass {
		pass[i] = charset[rand.Intn(len(charset))]
	}
	return string(pass)
}

func GenerateShortID() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, 6)
	for i := range id {
		id[i] = idCharset[rand.Intn(len(idCharset))]
	}
	return string(id)
}
