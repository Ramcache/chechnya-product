package repositories

import (
	"chechnya-product/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
)

type DashboardRepositoryInterface interface {
	GetDashboardData(ctx context.Context) (*models.DashboardData, error)
}

type DashboardRepository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetDashboardData(ctx context.Context) (*models.DashboardData, error) {
	var data models.DashboardData

	// Общее количество заказов
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&data.TotalOrders)
	if err != nil {
		return nil, err
	}

	// Общая выручка
	err = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(total), 0) FROM orders`).Scan(&data.TotalRevenue)
	if err != nil {
		return nil, err
	}

	// Топ продукты
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.id, p.name, SUM(oi.quantity) as sold
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		GROUP BY p.id
		ORDER BY sold DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.TopProduct
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Sold); err != nil {
			return nil, err
		}
		data.TopProducts = append(data.TopProducts, p)
	}

	// Продажи по дням
	rows, err = r.db.QueryContext(ctx, `
		SELECT DATE(created_at), COUNT(*), SUM(total)
		FROM orders
		WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.DailySales
		if err := rows.Scan(&d.Date, &d.Orders, &d.Revenue); err != nil {
			return nil, err
		}
		data.SalesByDay = append(data.SalesByDay, d)
	}

	return &data, nil
}
