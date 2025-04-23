package models

type TopProduct struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Sold      int    `json:"sold"`
}

type DailySales struct {
	Date    string  `json:"date"`
	Orders  int     `json:"orders"`
	Revenue float64 `json:"revenue"`
}

type DashboardData struct {
	TotalOrders  int          `json:"total_orders"`
	TotalRevenue float64      `json:"total_revenue"`
	TopProducts  []TopProduct `json:"top_products"`
	SalesByDay   []DailySales `json:"sales_by_day"`
}
