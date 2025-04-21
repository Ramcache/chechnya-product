// main.go
package main

import (
	"chechnya-product/config"
	_ "chechnya-product/docs"
	"chechnya-product/internal/db"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/logger"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/routes"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"chechnya-product/internal/ws"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	dbConn, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	hub := ws.NewHub(logger)
	go hub.Run()

	// --- Repositories ---
	var (
		userRepo     repositories.UserRepository     = repositories.NewUserRepo(dbConn)
		cartRepo     repositories.CartRepository     = repositories.NewCartRepo(dbConn)
		productRepo  repositories.ProductRepository  = repositories.NewProductRepo(dbConn)
		orderRepo    repositories.OrderRepository    = repositories.NewOrderRepo(dbConn)
		categoryRepo repositories.CategoryRepository = repositories.NewCategoryRepo(dbConn)
	)

	// --- JWT Manager ---
	var jwtManager utils.JWTManagerInterface = utils.NewJWTManager(cfg.JWTSecret, 72*time.Hour)

	// --- Services ---
	var (
		userService     services.UserServiceInterface     = services.NewUserService(userRepo, jwtManager)
		cartService     services.CartServiceInterface     = services.NewCartService(cartRepo, productRepo)
		productService  services.ProductServiceInterface  = services.NewProductService(productRepo)
		orderService    services.OrderServiceInterface    = services.NewOrderService(cartRepo, orderRepo, productRepo, hub)
		categoryService services.CategoryServiceInterface = services.NewCategoryService(categoryRepo)
	)

	// --- Handlers ---
	var (
		userHandler     handlers.UserHandlerInterface     = handlers.NewUserHandler(userService, logger)
		cartHandler     handlers.CartHandlerInterface     = handlers.NewCartHandler(cartService, logger)
		productHandler  handlers.ProductHandlerInterface  = handlers.NewProductHandler(productService, logger)
		orderHandler    handlers.OrderHandlerInterface    = handlers.NewOrderHandler(orderService, logger)
		categoryHandler handlers.CategoryHandlerInterface = handlers.NewCategoryHandler(categoryService, logger)
		logHandler      handlers.LogHandlerInterface      = handlers.NewLogHandler(logger, "logs/app.log")
	)

	router := mux.NewRouter()
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.HandleFunc("/ws/orders", hub.HandleConnections)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// --- Роуты вынесены ---
	routes.RegisterPublicRoutes(router, userHandler, productHandler, categoryHandler, cartHandler, orderHandler)
	routes.RegisterPrivateRoutes(router, userHandler, jwtManager)
	routes.RegisterAdminRoutes(router, productHandler, orderHandler, categoryHandler, logHandler, jwtManager)

	// --- CORS ---
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	logger.Sugar().Infow("Server is running", "port", cfg.Port, "env", cfg.Env)
	if err := http.ListenAndServe(":"+cfg.Port, corsMiddleware.Handler(router)); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
