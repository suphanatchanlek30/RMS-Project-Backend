package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type DashboardRepository struct {
	DB *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

func (r *DashboardRepository) GetSummary(ctx context.Context) (*models.DashboardSummaryData, error) {
	resp := &models.DashboardSummaryData{}

	// Get today's sales
	sales := 0.0
	err := r.DB.QueryRow(ctx, `
		SELECT COALESCE(SUM(p.total_amount), 0)
		FROM payments p
		WHERE p.payment_status = 'PAID'
		  AND DATE(p.payment_time) = CURRENT_DATE
	`).Scan(&sales)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	resp.TodaySales = sales

	// Get today's orders count
	orders := 0
	err = r.DB.QueryRow(ctx, `
		SELECT COALESCE(COUNT(co.order_id), 0)
		FROM customer_orders co
		WHERE DATE(co.order_time) = CURRENT_DATE
	`).Scan(&orders)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	resp.TodayOrders = orders

	// Get occupied tables
	occupied := 0
	err = r.DB.QueryRow(ctx, `
		SELECT COALESCE(COUNT(*), 0)
		FROM restaurant_tables
		WHERE table_status = 'OCCUPIED'
	`).Scan(&occupied)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	resp.OccupiedTables = occupied

	// Get available tables
	available := 0
	err = r.DB.QueryRow(ctx, `
		SELECT COALESCE(COUNT(*), 0)
		FROM restaurant_tables
		WHERE table_status = 'AVAILABLE'
	`).Scan(&available)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	resp.AvailableTables = available

	// Get top selling menu today
	topMenu := models.TopMenuInfo{}
	err = r.DB.QueryRow(ctx, `
		SELECT m.menu_id, m.menu_name, COALESCE(SUM(oi.quantity), 0) as total_sold
		FROM order_items oi
		JOIN menus m ON oi.menu_id = m.menu_id
		JOIN customer_orders co ON oi.order_id = co.order_id
		WHERE DATE(co.order_time) = CURRENT_DATE
		GROUP BY m.menu_id, m.menu_name
		ORDER BY total_sold DESC
		LIMIT 1
	`).Scan(&topMenu.MenuID, &topMenu.MenuName, &topMenu.TotalSold)
	if err != nil {
		// ถ้าไม่มีข้อมูลวันนี้ ให้ใส่ 0
		topMenu = models.TopMenuInfo{MenuID: 0, MenuName: "", TotalSold: 0}
	}
	resp.TopMenu = topMenu

	return resp, nil
}
