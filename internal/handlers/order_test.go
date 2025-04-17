package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PlaceOrder_Success(t *testing.T) {
	resetDB(t)
	// Вставка данных
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("insert cart: %v", err)
	}
	_, err = db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (1, 1, 2)`)
	if err != nil {
		t.Fatalf("insert cart_items: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/order", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func Test_PlaceOrder_Empty(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("insert cart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/order", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_PlaceOrder_NoStock(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("insert cart: %v", err)
	}
	// Банан — 0 на складе
	_, err = db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (1, 2, 1)`)
	if err != nil {
		t.Fatalf("insert cart_items: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/order", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_GetUserOrders(t *testing.T) {
	resetDB(t)

	// Подготовка: вставка заказа
	_, err := db.Exec(`INSERT INTO orders (user_id, total) VALUES ($1, $2)`, testUserID, 99.99)
	if err != nil {
		t.Fatalf("failed to insert order: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Проверяем структуру JSON-ответа
	type Order struct {
		ID        int     `json:"ID"`
		UserID    int     `json:"UserID"`
		Total     float64 `json:"Total"`
		CreatedAt string  `json:"CreatedAt"`
	}

	var orders []Order
	err = json.NewDecoder(rec.Body).Decode(&orders)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(orders) == 0 {
		t.Fatal("no orders returned")
	}

	if orders[0].Total != 99.99 {
		t.Errorf("expected total 99.99, got: %.2f", orders[0].Total)
	}
}

func Test_GetAllOrders(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO orders (user_id, total) VALUES ($1, $2)`, testUserID, 123.45)
	if err != nil {
		t.Fatalf("insert order: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/orders", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_ExportOrdersCSV(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO orders (user_id, total) VALUES ($1, $2)`, testUserID, 77.77)
	if err != nil {
		t.Fatalf("insert order: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/orders/export", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("Order ID")) {
		t.Errorf("expected CSV header, got: %s", rec.Body.String())
	}
}
