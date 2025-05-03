package utils

import (
	"chechnya-product/internal/models"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"regexp"
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

var phoneRegex = regexp.MustCompile(`^\+\d{10,15}$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidatePhone проверяет формат международного номера телефона
func ValidatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone format: expected format +71234567890")
	}
	return nil
}

// ValidateEmail проверяет формат email
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateIdentifier проверяет, что значение — валидный телефон или email
func ValidateIdentifier(id string) error {
	if phoneRegex.MatchString(id) {
		return nil
	}
	if emailRegex.MatchString(id) {
		return nil
	}
	return errors.New("invalid identifier: expected phone like +71234567890 or email like user@example.com")
}

func TodayDate() string {
	return time.Now().Format("2006-01-02")
}
