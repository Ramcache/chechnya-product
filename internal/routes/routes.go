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
	announcement handlers.AnnouncementHandlerInterface,
	review handlers.ReviewHandlerInterface,
	jwt utils.JWTManagerInterface,
) {
	public := r.PathPrefix("/api").Subrouter()

	public.Use(middleware.RateLimitMiddleware)

	public.Use(middleware.OptionalJWTMiddleware(jwt))

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
	public.HandleFunc("/cart/bulk", cart.AddBulkToCart).Methods(http.MethodPost)
	public.HandleFunc("/cart/{product_id}", cart.UpdateItem).Methods(http.MethodPut)
	public.HandleFunc("/cart/{product_id}", cart.DeleteItem).Methods(http.MethodDelete)

	// Заказы
	public.HandleFunc("/order", order.PlaceOrder).Methods(http.MethodPost)
	public.HandleFunc("/orders", order.GetUserOrders).Methods(http.MethodGet)
	public.HandleFunc("/order-reviews", order.GetAllReview).Methods(http.MethodGet)
	public.HandleFunc("/orders/{id}/review", order.LeaveReview).Methods(http.MethodPatch)
	public.HandleFunc("/orders/{id}/review", order.GetReview).Methods(http.MethodGet)
	public.HandleFunc("/orders/{id}/repeat", order.RepeatOrder).Methods(http.MethodPost)
	public.HandleFunc("/orders/history", order.GetOrderHistory).Methods(http.MethodGet)
	public.HandleFunc("/orders/{id}/status", order.UpdateStatus).Methods(http.MethodPatch)
	public.HandleFunc("/orders/{id}", order.GetOrderByID).Methods(http.MethodGet)

	// Объявления
	public.HandleFunc("/announcements", announcement.GetAll).Methods(http.MethodGet)
	public.HandleFunc("/announcements/{id}", announcement.GetByID).Methods(http.MethodGet)

	// Отзывы
	public.HandleFunc("/products/{id}/reviews", review.GetReviews).Methods(http.MethodGet)
	public.HandleFunc("/products/{id}/reviews", review.AddReview).Methods(http.MethodPost)
	public.HandleFunc("/products/{id}/reviews", review.UpdateReview).Methods(http.MethodPut)
	public.HandleFunc("/products/{id}/reviews", review.DeleteReview).Methods(http.MethodDelete)

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
	user handlers.UserHandlerInterface,
	product handlers.ProductHandlerInterface,
	order handlers.OrderHandlerInterface,
	category handlers.CategoryHandlerInterface,
	logs handlers.LogHandlerInterface,
	dashboard handlers.DashboardHandlerInterface,
	jwt utils.JWTManagerInterface,
	announcement handlers.AnnouncementHandlerInterface,
	adminInterface handlers.AdminInterface,
) {
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.JWTMiddleware(jwt))
	admin.Use(middleware.OnlyAdmin())

	admin.HandleFunc("/truncate", adminInterface.TruncateTableHandler).Methods(http.MethodPost)
	admin.HandleFunc("/truncate/all", adminInterface.TruncateAllTablesHandler).Methods(http.MethodPost)

	admin.HandleFunc("/users", user.CreateUserByPhone).Methods(http.MethodPost)

	// Управление товарами
	admin.HandleFunc("/products/{id}/upload-photo", product.UploadPhoto).Methods(http.MethodPost)
	admin.HandleFunc("/products", product.Add).Methods(http.MethodPost)
	admin.HandleFunc("/products/bulk", product.AddBulk).Methods(http.MethodPost)

	admin.HandleFunc("/products/{id}", product.Update).Methods(http.MethodPut)
	admin.HandleFunc("/products/{id}", product.Patch).Methods(http.MethodPatch)
	admin.HandleFunc("/products/{id}", product.Delete).Methods(http.MethodDelete)

	// Управление заказами
	admin.HandleFunc("/orders", order.GetAllOrders).Methods(http.MethodGet)
	admin.HandleFunc("/orders/export", order.ExportOrdersCSV).Methods(http.MethodGet)
	admin.HandleFunc("/orders/{id}", order.DeleteOrder).Methods(http.MethodDelete)

	// Управление категориями
	admin.HandleFunc("/categories", category.Create).Methods(http.MethodPost)
	admin.HandleFunc("/categories/bulk", category.CreateBulk).Methods(http.MethodPost)
	admin.HandleFunc("/categories/{id}", category.Update).Methods(http.MethodPut)
	admin.HandleFunc("/categories/{id}", category.Delete).Methods(http.MethodDelete)

	// Просмотр логов
	admin.HandleFunc("/logs", logs.GetLogs).Methods(http.MethodGet)

	// Дэшборд
	admin.HandleFunc("/dashboard", dashboard.GetDashboard).Methods(http.MethodGet)

	// Объявления
	admin.HandleFunc("/announcements", announcement.Create).Methods(http.MethodPost)
	admin.HandleFunc("/announcements/{id}", announcement.Update).Methods(http.MethodPut)
	admin.HandleFunc("/announcements/{id}", announcement.Delete).Methods(http.MethodDelete)
}
