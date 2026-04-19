package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type ReportRepository struct {
	DB *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) GetSalesReport(ctx context.Context, dateFrom string, dateTo string, groupBy string) ([]models.SalesReportItem, error) {
	// Validate dates
	fromDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DATE_FROM")
	}

	toDate, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DATE_TO")
	}

	// Validate that dateFrom <= dateTo
	if fromDate.After(toDate) {
		return nil, fmt.Errorf("INVALID_DATE_RANGE")
	}

	// Validate groupBy
	if groupBy != "day" && groupBy != "month" {
		return nil, fmt.Errorf("INVALID_GROUP_BY")
	}

	var dateFormat string
	if groupBy == "day" {
		dateFormat = "YYYY-MM-DD"
	} else {
		dateFormat = "YYYY-MM"
	}

	query := fmt.Sprintf(`
		SELECT 
			TO_CHAR(p.payment_time, '%s') as date,
			COALESCE(SUM(p.total_amount), 0) as total_sales,
			COALESCE(COUNT(DISTINCT co.order_id), 0) as total_orders
		FROM payments p
		LEFT JOIN customer_orders co ON co.order_id = p.order_id
		WHERE p.payment_status = 'PAID'
		  AND DATE(p.payment_time) >= $1
		  AND DATE(p.payment_time) <= $2
		GROUP BY TO_CHAR(p.payment_time, '%s')
		ORDER BY date ASC
	`, dateFormat, dateFormat)

	rows, err := r.DB.Query(ctx, query, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	var items []models.SalesReportItem
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

func (r *ReportRepository) GetTopMenusReport(ctx context.Context, dateFrom string, dateTo string, limit int) ([]models.TopMenuReportItem, error) {
	fromDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DATE_FROM")
	}

	toDate, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return nil, fmt.Errorf("INVALID_DATE_TO")
	}

	if fromDate.After(toDate) {
		return nil, fmt.Errorf("INVALID_DATE_RANGE")
	}

	if limit <= 0 {
		return nil, fmt.Errorf("INVALID_LIMIT")
	}

	query := `
		SELECT
			m.menu_id,
			m.menu_name,
			COALESCE(SUM(oi.quantity), 0) AS total_quantity,
			COALESCE(SUM(oi.quantity * oi.unit_price), 0) AS total_amount
		FROM order_items oi
		JOIN menus m ON oi.menu_id = m.menu_id
		JOIN customer_orders co ON oi.order_id = co.order_id
		WHERE DATE(co.order_time) >= $1
		  AND DATE(co.order_time) <= $2
		  AND oi.item_status <> 'CANCELLED'
		GROUP BY m.menu_id, m.menu_name
		ORDER BY total_quantity DESC
		LIMIT $3
	`

	rows, err := r.DB.Query(ctx, query, fromDate, toDate, limit)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	var items []models.TopMenuReportItem
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
