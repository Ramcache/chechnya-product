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
	"net/http"
)

func main() {
	// 📋 Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 🔧 Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	// 🛢️ Подключение к базе данных
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// 🧱 Инициализация репозиториев
	userRepo := repositories.NewUserRepo(database)
	cartRepo := repositories.NewCartRepo(database)
	productRepo := repositories.NewProductRepo(database)
	orderRepo := repositories.NewOrderRepo(database)

	// 🧠 Сервисы
	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo)

	// 🎮 Обработчики
	userHandler := handlers.NewUserHandler(userService)
	cartHandler := handlers.NewCartHandler(cartService, logger)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// 🚦 Роутер
	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware(logger))

	// 📚 Swagger-документация
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 🔓 Публичные маршруты
	public := router.PathPrefix("/api").Subrouter()
	public.HandleFunc("/register", userHandler.Register).Methods("POST")
	public.HandleFunc("/login", userHandler.Login).Methods("POST")
	public.HandleFunc("/products", productHandler.GetAll).Methods("GET")
	public.HandleFunc("/categories", productHandler.GetCategories).Methods("GET")
	public.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET")

	// 🔐 Приватные маршруты (для авторизованных пользователей)
	private := router.PathPrefix("/api").Subrouter()
	private.Use(middleware.JWTAuth(cfg.JWTSecret))

	// 📦 Корзина
	private.HandleFunc("/cart", cartHandler.AddToCart).Methods("POST")
	private.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")
	private.HandleFunc("/cart/{product_id}", cartHandler.UpdateItem).Methods("PUT")
	private.HandleFunc("/cart/{product_id}", cartHandler.DeleteItem).Methods("DELETE")
	private.HandleFunc("/cart/clear", cartHandler.ClearCart).Methods("DELETE")
	private.HandleFunc("/cart/checkout", cartHandler.Checkout).Methods("POST")

	// 🛒 Заказы
	private.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")
	private.HandleFunc("/orders", orderHandler.GetUserOrders).Methods("GET")

	// 👤 Профиль
	private.HandleFunc("/me", userHandler.Me).Methods("GET")

	// 🛠️ Админ-панель
	admin := router.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.JWTAuth(cfg.JWTSecret))
	admin.Use(middleware.OnlyAdmin())

	admin.HandleFunc("/products", productHandler.Add).Methods("POST")
	admin.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE")
	admin.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT")
	admin.HandleFunc("/orders", orderHandler.GetAllOrders).Methods("GET")
	admin.HandleFunc("/orders/export", orderHandler.ExportOrdersCSV).Methods("GET")

	// 🚀 Запуск сервера
	logger.Sugar().Infow("Server is running",
		"port", cfg.Port,
		"env", cfg.Env,
	)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
