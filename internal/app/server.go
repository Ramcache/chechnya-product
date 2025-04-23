package app

import (
	"chechnya-product/config"
	_ "chechnya-product/docs"
	"chechnya-product/internal/db"
	"chechnya-product/internal/handlers"
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
	"net/http"
	"time"
)

func NewServer(cfg *config.Config, logger *zap.Logger) *http.Server {
	dbConn, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	// 💡 ты можешь обернуть dbConn.Close() через defer в main, если захочешь.

	hub := ws.NewHub(logger)
	go hub.Run()

	// --- Repositories ---
	userRepo := repositories.NewUserRepo(dbConn)
	cartRepo := repositories.NewCartRepo(dbConn)
	productRepo := repositories.NewProductRepo(dbConn)
	orderRepo := repositories.NewOrderRepo(dbConn)
	categoryRepo := repositories.NewCategoryRepo(dbConn)
	dashboardRepo := repositories.NewDashboardRepository(dbConn)

	// --- JWT ---
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 72*time.Hour)

	// --- Services ---
	cartService := services.NewCartService(cartRepo, productRepo)
	userService := services.NewUserService(userRepo, jwtManager, cartService)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo, hub)
	categoryService := services.NewCategoryService(categoryRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// --- Handlers ---
	userHandler := handlers.NewUserHandler(userService, logger)
	cartHandler := handlers.NewCartHandler(cartService, logger)
	productHandler := handlers.NewProductHandler(productService, logger)
	orderHandler := handlers.NewOrderHandler(orderService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	logHandler := handlers.NewLogHandler(logger, "logs/app.log")
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, logger)

	// --- Router ---
	router := mux.NewRouter()
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.HandleFunc("/ws/orders", hub.HandleConnections)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	routes.RegisterPublicRoutes(router, userHandler, productHandler, categoryHandler, cartHandler, orderHandler)
	routes.RegisterPrivateRoutes(router, userHandler, jwtManager)
	routes.RegisterAdminRoutes(router, productHandler, orderHandler, categoryHandler, logHandler, dashboardHandler, jwtManager)

	// --- CORS ---
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	// --- HTTP Server ---
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: corsMiddleware.Handler(router),
	}
}
