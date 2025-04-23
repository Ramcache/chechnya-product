package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
