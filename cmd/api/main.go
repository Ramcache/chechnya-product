// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"chechnya-product/config"
	_ "chechnya-product/docs"
	"chechnya-product/internal/db"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/services"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
)

// @title        Chechnya Product API
// @version      1.0
// @description  Онлайн-магазин для продажи продуктов
// @host         localhost:8080
// @BasePath     /api
// @schemes      http
func main() {
	cfg := config.LoadConfig()

	// логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// БД
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД", zap.Error(err))
	}
	defer database.Close()

	// Роутер
	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware(logger))
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Репозитории и сервисы
	userRepo := repositories.NewUserRepo(database)
	cartRepo := repositories.NewCartRepo(database)
	productRepo := repositories.NewProductRepo(database)
	orderRepo := repositories.NewOrderRepo(database)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo)

	userHandler := handlers.NewUserHandler(userService)
	cartHandler := handlers.NewCartHandler(cartService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// 🔓 Публичные маршруты
	public := r.PathPrefix("/api").Subrouter()
	public.HandleFunc("/register", userHandler.Register).Methods("POST")
	public.HandleFunc("/login", userHandler.Login).Methods("POST")
	public.HandleFunc("/products", productHandler.GetAll).Methods("GET")
	public.HandleFunc("/categories", productHandler.GetCategories).Methods("GET")
	public.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET")

	// 🔐 Приватные маршруты (JWT Middleware)
	private := r.PathPrefix("/api").Subrouter()
	private.Use(middleware.JWTAuth(cfg.JWTSecret))
	private.HandleFunc("/cart", cartHandler.AddToCart).Methods("POST")
	private.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")
	private.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")
	private.HandleFunc("/orders", orderHandler.GetUserOrders).Methods("GET")
	private.HandleFunc("/me", userHandler.Me).Methods("GET")

	private.HandleFunc("/cart/{product_id}", cartHandler.UpdateItem).Methods("PUT")
	private.HandleFunc("/cart/{product_id}", cartHandler.DeleteItem).Methods("DELETE")

	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.JWTAuth(cfg.JWTSecret))
	admin.Use(middleware.OnlyAdmin())

	admin.HandleFunc("/products", productHandler.Add).Methods("POST")
	admin.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE")
	admin.HandleFunc("/orders", orderHandler.GetAllOrders).Methods("GET")
	admin.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT")
	admin.HandleFunc("/orders/export", orderHandler.ExportOrdersCSV).Methods("GET")

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
