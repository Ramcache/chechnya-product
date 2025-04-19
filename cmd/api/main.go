// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"chechnya-product/config"
	_ "chechnya-product/docs"
	"chechnya-product/internal/db"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/logger"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/services"
	"chechnya-product/internal/ws"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	//test
	// 📋 Инициализация логгера
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
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

	hub := ws.NewHub(logger)
	go hub.Run()

	// 🧱 Инициализация репозиториев
	userRepo := repositories.NewUserRepo(database)
	cartRepo := repositories.NewCartRepo(database)
	productRepo := repositories.NewProductRepo(database)
	orderRepo := repositories.NewOrderRepo(database)
	categoryRepo := repositories.NewCategoryRepo(database)
	verificationRepo := repositories.NewVerificationRepository(database)

	// 🧠 Сервисы
	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo, hub)
	categoryService := services.NewCategoryService(categoryRepo)
	verificationService := services.NewVerificationService(verificationRepo, "79298974969") // твой номер без +

	// 🎮 Обработчики
	userHandler := handlers.NewUserHandler(userService, logger)
	cartHandler := handlers.NewCartHandler(cartService, logger)
	productHandler := handlers.NewProductHandler(productService, logger)
	orderHandler := handlers.NewOrderHandler(orderService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	logHandler := handlers.NewLogHandler(logger, "logs/app.log")
	handler := handlers.NewVerificationHandler(verificationService)

	// 🚦 Роутер
	router := mux.NewRouter()

	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.HandleFunc("/ws/orders", hub.HandleConnections)

	router.HandleFunc("/verify/start", handler.StartVerification).Methods("POST")
	router.HandleFunc("/verify/confirm", handler.ConfirmCode).Methods("POST")

	// 📚 Swagger-документация
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// 🔓 Публичные маршруты
	public := router.PathPrefix("/api").Subrouter()
	public.HandleFunc("/register", userHandler.Register).Methods("POST")
	public.HandleFunc("/login", userHandler.Login).Methods("POST")
	public.HandleFunc("/products", productHandler.GetAll).Methods("GET")
	public.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET")
	public.HandleFunc("/categories", categoryHandler.GetAll).Methods("GET")

	// 📦 Корзина (доступна всем)
	public.HandleFunc("/cart", cartHandler.AddToCart).Methods("POST")
	public.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")
	public.HandleFunc("/cart/{product_id}", cartHandler.UpdateItem).Methods("PUT")
	public.HandleFunc("/cart/{product_id}", cartHandler.DeleteItem).Methods("DELETE")
	public.HandleFunc("/cart/clear", cartHandler.ClearCart).Methods("DELETE")
	public.HandleFunc("/cart/checkout", cartHandler.Checkout).Methods("POST")

	// 🛒 Заказы (доступны всем)
	public.HandleFunc("/order", orderHandler.PlaceOrder).Methods("POST")
	public.HandleFunc("/orders", orderHandler.GetUserOrders).Methods("GET")

	// 🔐 Приватные маршруты
	private := router.PathPrefix("/api").Subrouter()
	private.Use(middleware.JWTAuth(cfg.JWTSecret))
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
	admin.HandleFunc("/categories", categoryHandler.Create).Methods("POST")
	admin.HandleFunc("/categories/{id}", categoryHandler.Update).Methods("PUT")
	admin.HandleFunc("/categories/{id}", categoryHandler.Delete).Methods("DELETE")
	admin.HandleFunc("/logs", logHandler.GetLogs).Methods("GET")

	// 🛡️ CORS Middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	// 🚀 Запуск сервера с CORS
	logger.Sugar().Infow("Server is running", "port", cfg.Port, "env", cfg.Env)

	if err := http.ListenAndServe(":"+cfg.Port, corsMiddleware.Handler(router)); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
