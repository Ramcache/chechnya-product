package handlers_test

import (
	"bytes"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/services"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sqlx.DB
var router *mux.Router

const testUserID = 1

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=myshop password=625325 sslmode=disable")
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	productRepo := repositories.NewProductRepo(db)
	cartRepo := repositories.NewCartRepo(db)
	service := services.NewCartService(cartRepo, productRepo)
	logger := zap.NewNop()
	handler := handlers.NewCartHandler(service, logger)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	orderRepo := repositories.NewOrderRepo(db)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	router = mux.NewRouter()
	router.Use(fakeUserMiddleware)
	router.HandleFunc("/cart", handler.AddToCart).Methods("POST")
	router.HandleFunc("/cart", handler.GetCart).Methods("GET")
	router.HandleFunc("/cart/{product_id}", handler.UpdateItem).Methods("PUT")
	router.HandleFunc("/cart/{product_id}", handler.DeleteItem).Methods("DELETE")
	router.HandleFunc("/cart/clear", handler.ClearCart).Methods("DELETE")
	router.HandleFunc("/cart/checkout", handler.Checkout).Methods("POST")

	router.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")
	router.HandleFunc("/orders", orderHandler.GetUserOrders).Methods("GET")
	router.HandleFunc("/admin/orders", orderHandler.GetAllOrders).Methods("GET")
	router.HandleFunc("/admin/orders/export", orderHandler.ExportOrdersCSV).Methods("GET")

	router.HandleFunc("/products", productHandler.GetAll).Methods("GET")
	router.HandleFunc("/products/categories", productHandler.GetCategories).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET")
	router.HandleFunc("/products", productHandler.Add).Methods("POST")
	router.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT")
	router.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE")

	fmt.Println("Connected to DB:", db.DriverName())
	var dbName string
	_ = db.Get(&dbName, "SELECT current_database()")
	fmt.Println("Current database:", dbName)

	os.Exit(m.Run())
}

func fakeUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.ContextUserID, testUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func resetDB(t *testing.T) {
	t.Helper()
	_, err := db.Exec(`
		TRUNCATE cart_items, carts, products, users, orders RESTART IDENTITY CASCADE;
		INSERT INTO users (id, username, password, role) VALUES 
			(1, 'testuser', 'hashed-password', 'user'),
			(2, 'adminuser', 'hashed-password', 'admin');

		INSERT INTO products (name, price, stock, description, category) 
		VALUES 
			('Apple', 10.0, 5, 'Red apple', 'fruit'), 
			('Banana', 5.0, 0, 'Ripe banana', 'fruit');
	`)
	if err != nil {
		t.Fatalf("failed to reset db: %v", err)
	}
}

func Test_AddToCart_Success(t *testing.T) {
	resetDB(t)

	// 🔧 Явно создаём корзину для testUserID (иначе FOREIGN KEY падает)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("failed to insert cart: %v", err)
	}

	body := `{"product_id": 1, "quantity": 2}`
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Logf("response body: %s", rec.Body.String())
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func Test_AddToCart_ProductNotFound(t *testing.T) {
	resetDB(t)
	body := `{"product_id": 999, "quantity": 1}`
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_AddToCart_StockExceeded(t *testing.T) {
	resetDB(t)
	body := `{"product_id": 1, "quantity": 10}`
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_AddToCart_InvalidJSON(t *testing.T) {
	resetDB(t)
	body := `{"product_id": "apple", "quantity": 2}`
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_AddToCart_ZeroQuantity(t *testing.T) {
	resetDB(t)
	body := `{"product_id": 1, "quantity": 0}`
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_AddToCart_NegativeQuantity(t *testing.T) {
	resetDB(t)
	body := `{"product_id": 1, "quantity": -5}`
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_AddToCart_DuplicateProduct(t *testing.T) {
	resetDB(t)

	// 🔧 Вставляем корзину вручную, чтобы избежать ошибки FK
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("failed to insert cart: %v", err)
	}

	body := `{"product_id": 1, "quantity": 2}`

	// Первый запрос
	req1 := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusCreated {
		t.Logf("response body: %s", rec1.Body.String())
		t.Fatalf("expected 201 on first add, got %d", rec1.Code)
	}

	// Второй запрос — новый req!
	req2 := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(body))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusCreated {
		t.Logf("response body: %s", rec2.Body.String())
		t.Fatalf("expected 201 on duplicate add, got %d", rec2.Code)
	}
}

func Test_GetCart_ReturnsCorrectData(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("failed to insert cart: %v", err)
	}
	_, err = db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (1, 1, 2)`)
	if err != nil {
		t.Fatalf("failed to insert cart item: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/cart", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte(`"product_id":1`)) {
		t.Errorf("expected product_id in response, got %s", rec.Body.String())
	}
}

func Test_UpdateItem_Success(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("failed to insert cart: %v", err)
	}
	_, err = db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (1, 1, 2)`)
	if err != nil {
		t.Fatalf("failed to insert cart item: %v", err)
	}

	body := `{"quantity": 3}`
	req := httptest.NewRequest(http.MethodPut, "/cart/1", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_UpdateItem_TooMuch(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("failed to insert cart: %v", err)
	}
	_, err = db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (1, 1, 2)`)
	if err != nil {
		t.Fatalf("failed to insert cart item: %v", err)
	}

	body := `{"quantity": 999}`
	req := httptest.NewRequest(http.MethodPut, "/cart/1", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_DeleteItem_Success(t *testing.T) {
	resetDB(t)
	_, err := db.Exec(`INSERT INTO carts (user_id) VALUES ($1)`, testUserID)
	if err != nil {
		t.Fatalf("failed to insert cart: %v", err)
	}
	_, err = db.Exec(`INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (1, 1, 2)`)
	if err != nil {
		t.Fatalf("failed to insert cart item: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/cart/1", nil)
	req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_ClearCart(t *testing.T) {
	resetDB(t)
	req := httptest.NewRequest(http.MethodDelete, "/cart/clear", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func Test_Checkout(t *testing.T) {
	resetDB(t)
	req := httptest.NewRequest(http.MethodPost, "/cart/checkout", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
