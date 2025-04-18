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
	"chechnya-product/internal/ws"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Получение идентификатора пользователя или IP
func getUserIdentifier(r *http.Request) string {
	userID := middleware.GetUserID(r)
	if userID != 0 {
		return fmt.Sprintf("user_%d", userID)
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return fmt.Sprintf("ip_%s", ip)
}

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

	hub := ws.NewHub()
	go hub.Run()

	// 🧱 Инициализация репозиториев
	userRepo := repositories.NewUserRepo(database)
	cartRepo := repositories.NewCartRepo(database)
	productRepo := repositories.NewProductRepo(database)
	orderRepo := repositories.NewOrderRepo(database)

	// 🧠 Сервисы
	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo, hub)

	// 🎮 Обработчики
	userHandler := handlers.NewUserHandler(userService)
	cartHandler := handlers.NewCartHandler(cartService, logger)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	categoryRepo := repositories.NewCategoryRepo(database)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// 🚦 Роутер
	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware(logger))
	router.HandleFunc("/ws/orders", hub.HandleConnections)

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
	admin.HandleFunc("/categories", categoryHandler.Create).Methods("POST")
	admin.HandleFunc("/categories/{id}", categoryHandler.Update).Methods("PUT")
	admin.HandleFunc("/categories/{id}", categoryHandler.Delete).Methods("DELETE")

	// 🚀 Запуск сервера
	logger.Sugar().Infow("Server is running", "port", cfg.Port, "env", cfg.Env)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
