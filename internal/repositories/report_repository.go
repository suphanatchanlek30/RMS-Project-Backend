package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type ReportRepository struct {
	DB *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) GetSalesReport(ctx context.Context, query models.SalesReportQuery) ([]models.SalesReportItem, error) {
	var dateExpr string
	switch query.GroupBy {
	case "day":
		dateExpr = "TO_CHAR(p.payment_time AT TIME ZONE 'UTC', 'YYYY-MM-DD')"
	case "month":
		dateExpr = "TO_CHAR(p.payment_time AT TIME ZONE 'UTC', 'YYYY-MM')"
	default:
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	sql := fmt.Sprintf(`
		SELECT
			period,
			COALESCE(SUM(total_amount), 0) AS total_sales,
			COALESCE(SUM(order_count), 0)::int AS total_orders
		FROM (
			SELECT
				%s AS period,
				p.total_amount,
				(
					SELECT COUNT(*)
					FROM customer_orders co
					WHERE co.session_id = paid_order.session_id
				) AS order_count
			FROM payments p
			JOIN customer_orders paid_order ON paid_order.order_id = p.order_id
			WHERE p.payment_status = 'PAID'
			  AND p.payment_time >= $1
			  AND p.payment_time <= $2
		) report_rows
		GROUP BY period
		ORDER BY period ASC
	`, dateExpr)

	rows, err := r.DB.Query(ctx, sql, query.DateFrom, query.DateTo)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	items := make([]models.SalesReportItem, 0)
	for rows.Next() {
		var item models.SalesReportItem
		if err := rows.Scan(&item.Date, &item.TotalSales, &item.TotalOrders); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return items, nil
}

func (r *ReportRepository) GetTopMenusReport(ctx context.Context, query models.TopMenusReportQuery) ([]models.TopMenuReportItem, error) {
	sql := `
		SELECT
			m.menu_id,
			m.menu_name,
			SUM(oi.quantity)::int AS total_quantity,
			COALESCE(SUM(oi.quantity * oi.unit_price), 0) AS total_amount
		FROM payments p
		JOIN customer_orders paid_order ON paid_order.order_id = p.order_id
		JOIN customer_orders co ON co.session_id = paid_order.session_id
		JOIN order_items oi ON oi.order_id = co.order_id
		JOIN menus m ON m.menu_id = oi.menu_id
		WHERE p.payment_status = 'PAID'
		  AND p.payment_time >= $1
		  AND p.payment_time <= $2
		  AND oi.item_status <> 'CANCELLED'
		GROUP BY m.menu_id, m.menu_name
		ORDER BY total_quantity DESC, m.menu_id ASC
		LIMIT $3
	`

	rows, err := r.DB.Query(ctx, sql, query.DateFrom, query.DateTo, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	items := make([]models.TopMenuReportItem, 0)
	for rows.Next() {
		var item models.TopMenuReportItem
		if err := rows.Scan(&item.MenuID, &item.MenuName, &item.TotalQuantity, &item.TotalAmount); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return items, nil
}
