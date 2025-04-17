package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetAllProducts(t *testing.T) {
	resetDB(t)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_GetProductByID_Success(t *testing.T) {
	resetDB(t)
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_GetProductByID_NotFound(t *testing.T) {
	resetDB(t)
	req := httptest.NewRequest(http.MethodGet, "/products/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func Test_AddProduct_AsAdmin(t *testing.T) {
	resetDB(t)

	body := `{
		"name": "Test Product",
		"description": "A test product",
		"price": 12.5,
		"stock": 10,
		"category": "test-category"
	}`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), "role", "admin")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("body: %s", rec.Body.String())
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func Test_AddProduct_Forbidden(t *testing.T) {
	resetDB(t)

	body := `{"name":"Test","description":"Desc","price":12.5,"stock":5}`
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "role", "user") // 🛑 нет доступа
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func Test_UpdateProduct(t *testing.T) {
	resetDB(t)

	body := `{
		"name": "Updated Apple",
		"description": "Green apple",
		"price": 15.0,
		"stock": 7,
		"category": "fruits"
	}`

	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	ctx := context.WithValue(req.Context(), "role", "admin")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("body: %s", rec.Body.String())
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_DeleteProduct_AsAdmin(t *testing.T) {
	resetDB(t)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), "role", "admin")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_GetCategories(t *testing.T) {
	resetDB(t)

	_, err := db.Exec(`UPDATE products SET category = 'fruit' WHERE id = 1`)
	if err != nil {
		t.Fatalf("failed to update category: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/products/categories", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var categories []string
	if err := json.NewDecoder(rec.Body).Decode(&categories); err != nil {
		t.Fatalf("failed to decode categories: %v", err)
	}
	if len(categories) == 0 {
		t.Error("expected at least 1 category")
	}
}
