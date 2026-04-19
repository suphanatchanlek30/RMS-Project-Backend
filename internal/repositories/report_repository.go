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
