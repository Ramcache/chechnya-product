package app

import (
	"chechnya-product/config"
	_ "chechnya-product/docs"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/routes"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"chechnya-product/internal/ws"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewServer(cfg *config.Config, logger *zap.Logger, dbConn *sqlx.DB) *http.Server {
	hub := ws.NewHub(logger)
	go hub.Run()

	// --- Repositories ---
	userRepo := repositories.NewUserRepo(dbConn)
	cartRepo := repositories.NewCartRepo(dbConn)
	productRepo := repositories.NewProductRepo(dbConn)
	orderRepo := repositories.NewOrderRepo(dbConn)
	categoryRepo := repositories.NewCategoryRepo(dbConn)
	dashboardRepo := repositories.NewDashboardRepository(dbConn)
	announcementRepo := repositories.NewAnnouncementRepo(dbConn)
	reviewRepo := repositories.NewReviewRepo(dbConn)
	// --- JWT ---
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 72*time.Hour)

	// --- Services ---
	cartService := services.NewCartService(cartRepo, productRepo)
	userService := services.NewUserService(userRepo, jwtManager, cartService)
	productService := services.NewProductService(productRepo, logger)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo, hub)
	categoryService := services.NewCategoryService(categoryRepo, logger)
	dashboardService := services.NewDashboardService(dashboardRepo)
	announcementService := services.NewAnnouncementService(announcementRepo, hub)
	reviewService := services.NewReviewService(reviewRepo)

	// --- Handlers ---
	userHandler := handlers.NewUserHandler(userService, logger)
	cartHandler := handlers.NewCartHandler(cartService, logger)
	productHandler := handlers.NewProductHandler(productService, logger)
	orderHandler := handlers.NewOrderHandler(orderService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	logHandler := handlers.NewLogHandler(logger, "logs/app.log")
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, logger)
	announcementHandler := handlers.NewAnnouncementHandler(announcementService, logger)
	reviewHandler := handlers.NewReviewHandler(reviewService, logger)

	// --- Router ---
	router := mux.NewRouter()
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.HandleFunc("/ws/orders", hub.HandleConnections)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	routes.RegisterPublicRoutes(router, userHandler, productHandler, categoryHandler, cartHandler, orderHandler, announcementHandler, reviewHandler, jwtManager)
	routes.RegisterPrivateRoutes(router, userHandler, jwtManager)
	routes.RegisterAdminRoutes(router, productHandler, orderHandler, categoryHandler, logHandler, dashboardHandler, jwtManager, announcementHandler)

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
