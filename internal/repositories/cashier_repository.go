package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type CashierRepository struct {
	DB *pgxpool.Pool
}

func NewCashierRepository(db *pgxpool.Pool) *CashierRepository {
	return &CashierRepository{DB: db}
}

func (r *CashierRepository) GetTablesOverview(ctx context.Context) ([]models.CashierTableOverview, error) {
	query := `
		SELECT
			rt.table_id,
			rt.table_number,
			rt.table_status,
			ts.session_id,
			ts.start_time
		FROM restaurant_tables rt
		LEFT JOIN table_sessions ts
			ON ts.table_id = rt.table_id
			AND ts.session_status = 'OPEN'
		ORDER BY rt.table_number ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.CashierTableOverview

	for rows.Next() {
		var t models.CashierTableOverview
		var sessionID *int
		var startTime *time.Time

		if err := rows.Scan(
			&t.TableID,
			&t.TableNumber,
			&t.TableStatus,
			&sessionID,
			&startTime,
		); err != nil {
			return nil, err
		}

		if sessionID != nil && startTime != nil {
			t.CurrentSession = &models.CashierCurrentSession{
				SessionID: *sessionID,
				StartTime: *startTime,
			}
		} else {
			t.CurrentSession = nil
		}

		tables = append(tables, t)
	}

	return tables, rows.Err()
}

func (r *CashierRepository) GetCheckout(ctx context.Context, sessionID int) (*models.CheckoutResponse, error) {
	query := `SELECT ts.session_id, rt.table_id, rt.table_number
	FROM table_sessions ts
	JOIN restaurant_tables rt ON ts.table_id = rt.table_id
	WHERE ts.session_id = $1`

	var resp models.CheckoutResponse
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(&resp.SessionID, &resp.TableID, &resp.TableNumber)
	if err != nil {
		return nil, err
	}

	query = `SELECT oi.order_item_id, m.menu_name, oi.quantity, oi.unit_price, (oi.quantity * oi.unit_price) as line_total
	FROM order_items oi
	JOIN menus m ON oi.menu_id = m.menu_id
	JOIN customer_orders co ON oi.order_id = co.order_id
	WHERE co.session_id = $1`

	rows, err := r.DB.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total float64
	for rows.Next() {
		var item models.CheckoutBillItem
		err := rows.Scan(&item.OrderItemID, &item.MenuName, &item.Quantity, &item.UnitPrice, &item.LineTotal)
		if err != nil {
			return nil, err
		}
		resp.Bill.Items = append(resp.Bill.Items, item)
		total += item.LineTotal
	}

	resp.Bill.TotalAmount = total

	query = `SELECT payment_method_id, method_name FROM payment_methods`

	rows2, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var pm models.CheckoutPaymentMethod
		err := rows2.Scan(&pm.PaymentMethodID, &pm.MethodName)
		if err != nil {
			return nil, err
		}
		resp.PaymentMethods = append(resp.PaymentMethods, pm)
	}

	return &resp, nil
}