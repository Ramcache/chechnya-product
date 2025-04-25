package routes

import (
	"chechnya-product/internal/handlers"
	"chechnya-product/internal/middleware"
	"chechnya-product/internal/utils"

	"github.com/gorilla/mux"
	"net/http"
)

func RegisterPublicRoutes(
	r *mux.Router,
	user handlers.UserHandlerInterface,
	product handlers.ProductHandlerInterface,
	category handlers.CategoryHandlerInterface,
	cart handlers.CartHandlerInterface,
	order handlers.OrderHandlerInterface,
) {
	public := r.PathPrefix("/api").Subrouter()

	// Аутентификация и регистрация
	public.HandleFunc("/register", user.Register).Methods(http.MethodPost)
	public.HandleFunc("/login", user.Login).Methods(http.MethodPost)

	// Товары и категории
	public.HandleFunc("/products", product.GetAll).Methods(http.MethodGet)
	public.HandleFunc("/products/{id}", product.GetByID).Methods(http.MethodGet)
	public.HandleFunc("/categories", category.GetAll).Methods(http.MethodGet)

	// Корзина
	public.HandleFunc("/cart", cart.AddToCart).Methods(http.MethodPost)
	public.HandleFunc("/cart", cart.GetCart).Methods(http.MethodGet)
	public.HandleFunc("/cart/clear", cart.ClearCart).Methods(http.MethodDelete)
	public.HandleFunc("/cart/{product_id}", cart.UpdateItem).Methods(http.MethodPut)
	public.HandleFunc("/cart/{product_id}", cart.DeleteItem).Methods(http.MethodDelete)

	// Заказы
	public.HandleFunc("/order", order.PlaceOrder).Methods(http.MethodPost)
	public.HandleFunc("/orders", order.GetUserOrders).Methods(http.MethodGet)

}

func RegisterPrivateRoutes(
	r *mux.Router,
	user handlers.UserHandlerInterface,
	jwt utils.JWTManagerInterface,
) {
	private := r.PathPrefix("/api").Subrouter()
	private.Use(middleware.JWTMiddleware(jwt))
	private.HandleFunc("/me", user.Me).Methods(http.MethodGet)
}

func RegisterAdminRoutes(
	r *mux.Router,
	product handlers.ProductHandlerInterface,
	order handlers.OrderHandlerInterface,
	category handlers.CategoryHandlerInterface,
	logs handlers.LogHandlerInterface,
	dashboard handlers.DashboardHandlerInterface,
	jwt utils.JWTManagerInterface,
) {
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.JWTMiddleware(jwt))
	admin.Use(middleware.OnlyAdmin())

	// Управление товарами
	admin.HandleFunc("/products", product.Add).Methods(http.MethodPost)
	admin.HandleFunc("/products/{id}", product.Update).Methods(http.MethodPut)
	admin.HandleFunc("/products/{id}", product.Delete).Methods(http.MethodDelete)

	// Управление заказами
	admin.HandleFunc("/orders", order.GetAllOrders).Methods(http.MethodGet)
	admin.HandleFunc("/orders/export", order.ExportOrdersCSV).Methods(http.MethodGet)

	// Управление категориями
	admin.HandleFunc("/categories", category.Create).Methods(http.MethodPost)
	admin.HandleFunc("/categories/bulk", category.CreateBulk).Methods(http.MethodPost)
	admin.HandleFunc("/categories/{id}", category.Update).Methods(http.MethodPut)
	admin.HandleFunc("/categories/{id}", category.Delete).Methods(http.MethodDelete)

	// Просмотр логов
	admin.HandleFunc("/logs", logs.GetLogs).Methods(http.MethodGet)

	// Дэшборд
	admin.HandleFunc("/dashboard", dashboard.GetDashboard).Methods(http.MethodGet)
}
