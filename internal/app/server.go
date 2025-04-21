package app

import (
	"chechnya-product/config"
	"chechnya-product/internal/db"
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/repositories"
	"chechnya-product/internal/routes"
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"chechnya-product/internal/ws"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewServer(cfg *config.Config, logger *zap.Logger, redis *redis.Client) *http.Server {
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
	verificationRepo := repositories.NewVerificationRepository(dbConn)

	// --- JWT ---
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, 72*time.Hour)

	// --- Services ---
	cartService := services.NewCartService(cartRepo, productRepo)
	userService := services.NewUserService(userRepo, jwtManager, cartService)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(cartRepo, orderRepo, productRepo, hub)
	categoryService := services.NewCategoryService(categoryRepo)
	verificationService := services.NewVerificationService(verificationRepo, "79298974969")

	// --- Handlers ---
	userHandler := handlers.NewUserHandler(userService, logger)
	cartHandler := handlers.NewCartHandler(cartService, logger)
	productHandler := handlers.NewProductHandler(productService, logger)
	orderHandler := handlers.NewOrderHandler(orderService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	logHandler := handlers.NewLogHandler(logger, "logs/app.log")
	verificationHandler := handlers.NewVerificationHandler(verificationService, logger)

	// --- Router ---
	router := mux.NewRouter()
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.HandleFunc("/ws/orders", hub.HandleConnections)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/verify/start", verificationHandler.StartVerification).Methods("POST")
	router.HandleFunc("/verify/confirm", verificationHandler.ConfirmCode).Methods("POST")
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

	// --- HTTP Server ---
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: corsMiddleware.Handler(router),
	}
}
