package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type DashboardRepository struct {
	DB *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

func (r *DashboardRepository) GetSummary(ctx context.Context) (*models.DashboardSummaryResponse, error) {
	resp := &models.DashboardSummaryResponse{}

	salesQuery := `
		SELECT COALESCE(SUM(p.total_amount), 0)
		FROM payments p
		WHERE p.payment_status = 'PAID'
		  AND p.payment_time::date = CURRENT_DATE
	`
	if err := r.DB.QueryRow(ctx, salesQuery).Scan(&resp.TodaySales); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	ordersQuery := `
		SELECT COUNT(*)
		FROM customer_orders
		WHERE order_time::date = CURRENT_DATE
	`
	if err := r.DB.QueryRow(ctx, ordersQuery).Scan(&resp.TodayOrders); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if err := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM restaurant_tables WHERE table_status = 'OCCUPIED'`).Scan(&resp.OccupiedTables); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if err := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM restaurant_tables WHERE table_status = 'AVAILABLE'`).Scan(&resp.AvailableTables); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	topMenuQuery := `
		SELECT m.menu_id, m.menu_name, SUM(oi.quantity)::int AS total_sold
		FROM payments p
		JOIN customer_orders paid_order ON paid_order.order_id = p.order_id
		JOIN customer_orders co ON co.session_id = paid_order.session_id
		JOIN order_items oi ON oi.order_id = co.order_id
		JOIN menus m ON m.menu_id = oi.menu_id
		WHERE p.payment_status = 'PAID'
		  AND p.payment_time::date = CURRENT_DATE
		  AND oi.item_status <> 'CANCELLED'
		GROUP BY m.menu_id, m.menu_name
		ORDER BY total_sold DESC, m.menu_id ASC
		LIMIT 1
	`
	var topMenu models.TopMenuSummary
	if err := r.DB.QueryRow(ctx, topMenuQuery).Scan(&topMenu.MenuID, &topMenu.MenuName, &topMenu.TotalSold); err == nil {
		resp.TopMenu = &topMenu
	} else if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}
