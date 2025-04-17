package handlers_test

import (
	"bytes"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/services"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupUserTest(t *testing.T) *handlers.UserHandler {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("PORT", "8080") // если тоже требуется

	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=myshop password=625325 sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	repo := repositories.NewUserRepo(db)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	_, err = db.Exec("TRUNCATE users RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("failed to truncate users: %v", err)
	}

	return handler
}

func Test_Register_Success(t *testing.T) {
	handler := setupUserTest(t)

	reqBody := `{"username": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func Test_Register_InvalidData(t *testing.T) {
	handler := setupUserTest(t)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`invalid json`))
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func Test_Register_UsernameTaken(t *testing.T) {
	handler := setupUserTest(t)

	reqBody := `{"username": "taken", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	// повторно
	req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	handler.Register(rec2, req2)

	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for duplicate, got %d", rec2.Code)
	}
}

func Test_Login_Success(t *testing.T) {
	handler := setupUserTest(t)

	// Регистрация
	reqBody := `{"username": "loginuser", "password": "securepass"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.Register(rec, req)

	// Логин
	req2 := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	handler.Login(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}

	var resp map[string]string
	_ = json.NewDecoder(rec2.Body).Decode(&resp)
	if resp["token"] == "" {
		t.Error("token was not returned")
	}
}

func Test_Login_Failure(t *testing.T) {
	handler := setupUserTest(t)

	reqBody := `{"username": "wrong", "password": "wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
	rec := httptest.NewRecorder()
	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
