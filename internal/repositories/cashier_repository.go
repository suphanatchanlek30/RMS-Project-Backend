package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type CashierRepository struct {
	DB *pgxpool.Pool
}

func NewCashierRepository(db *pgxpool.Pool) *CashierRepository {
	return &CashierRepository{DB: db}
}

func (r *CashierRepository) GetTablesOverview(ctx context.Context) ([]models.CashierTableOverviewItem, error) {
	query := `
		SELECT
			rt.table_id,
			rt.table_number,
			rt.table_status,
			cs.session_id,
			cs.start_time
		FROM restaurant_tables rt
		LEFT JOIN LATERAL (
			SELECT ts.session_id, ts.start_time
			FROM table_sessions ts
			WHERE ts.table_id = rt.table_id
			  AND ts.session_status = 'OPEN'
			ORDER BY ts.start_time DESC
			LIMIT 1
		) cs ON TRUE
		ORDER BY rt.table_number ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	items := make([]models.CashierTableOverviewItem, 0)
	for rows.Next() {
		var item models.CashierTableOverviewItem
		var sessionID *int
		var startTime *time.Time
		if err := rows.Scan(&item.TableID, &item.TableNumber, &item.TableStatus, &sessionID, &startTime); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		if sessionID != nil && startTime != nil {
			item.CurrentSession = &models.CashierTableCurrentSession{
				SessionID: *sessionID,
				StartTime: *startTime,
			}
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return items, nil
}

func (r *CashierRepository) GetSessionCheckout(ctx context.Context, sessionID int) (*models.SessionCheckoutResponse, error) {
	infoQuery := `
		SELECT ts.session_id, ts.table_id, rt.table_number
		FROM table_sessions ts
		JOIN restaurant_tables rt ON rt.table_id = ts.table_id
		WHERE ts.session_id = $1
	`

	resp := &models.SessionCheckoutResponse{}
	if err := r.DB.QueryRow(ctx, infoQuery, sessionID).Scan(&resp.SessionID, &resp.TableID, &resp.TableNumber); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	itemsQuery := `
		SELECT
			oi.order_item_id,
			m.menu_name,
			oi.quantity,
			oi.unit_price,
			(oi.quantity * oi.unit_price) AS line_total
		FROM customer_orders co
		JOIN order_items oi ON oi.order_id = co.order_id
		JOIN menus m ON m.menu_id = oi.menu_id
		WHERE co.session_id = $1
		  AND oi.item_status <> 'CANCELLED'
		ORDER BY oi.order_item_id ASC
	`

	itemRows, err := r.DB.Query(ctx, itemsQuery, sessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer itemRows.Close()

	billItems := make([]models.BillItemResponse, 0)
	var totalAmount float64
	for itemRows.Next() {
		var item models.BillItemResponse
		if err := itemRows.Scan(&item.OrderItemID, &item.MenuName, &item.Quantity, &item.UnitPrice, &item.LineTotal); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		totalAmount += item.LineTotal
		billItems = append(billItems, item)
	}
	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	paymentQuery := `
		SELECT payment_method_id, method_name
		FROM payment_methods
		ORDER BY payment_method_id ASC
	`
	methodRows, err := r.DB.Query(ctx, paymentQuery)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer methodRows.Close()

	methods := make([]models.PaymentMethodItem, 0)
	for methodRows.Next() {
		var method models.PaymentMethodItem
		if err := methodRows.Scan(&method.PaymentMethodID, &method.MethodName); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		methods = append(methods, method)
	}
	if err := methodRows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	resp.Bill = models.CheckoutBill{
		Items:       billItems,
		TotalAmount: totalAmount,
	}
	resp.PaymentMethods = methods

	return resp, nil
}

func (r *CashierRepository) ProcessCheckout(ctx context.Context, req models.CashierCheckoutRequest) (*models.CashierCheckoutResult, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer tx.Rollback(ctx)

	var tableID int
	var sessionStatus string
	var tableStatus string
	sessionQuery := `
		SELECT ts.table_id, ts.session_status, rt.table_status
		FROM table_sessions ts
		JOIN restaurant_tables rt ON rt.table_id = ts.table_id
		WHERE ts.session_id = $1
		FOR UPDATE OF ts, rt
	`
	if err := tx.QueryRow(ctx, sessionQuery, req.SessionID).Scan(&tableID, &sessionStatus, &tableStatus); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("NOT_FOUND_SESSION")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	if sessionStatus == "CLOSED" || tableStatus != "OCCUPIED" {
		return nil, fmt.Errorf("CONFLICT")
	}

	var paymentMethodName string
	if err := tx.QueryRow(ctx, `SELECT method_name FROM payment_methods WHERE payment_method_id = $1`, req.PaymentMethodID).Scan(&paymentMethodName); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("NOT_FOUND_PAYMENT_METHOD")
		}
		return nil, fmt.Errorf("INTERNAL")
	}
	_ = paymentMethodName

	var hasPaid bool
	paidQuery := `
		SELECT EXISTS(
			SELECT 1
			FROM payments p
			JOIN customer_orders paid_order ON paid_order.order_id = p.order_id
			JOIN customer_orders co ON co.session_id = paid_order.session_id
			WHERE co.session_id = $1
			  AND p.payment_status = 'PAID'
		)
	`
	if err := tx.QueryRow(ctx, paidQuery, req.SessionID).Scan(&hasPaid); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasPaid {
		return nil, fmt.Errorf("CONFLICT")
	}

	var hasPendingItems bool
	pendingQuery := `
		SELECT EXISTS(
			SELECT 1
			FROM customer_orders co
			JOIN order_items oi ON oi.order_id = co.order_id
			WHERE co.session_id = $1
			  AND oi.item_status IN ('WAITING', 'PREPARING')
		)
	`
	if err := tx.QueryRow(ctx, pendingQuery, req.SessionID).Scan(&hasPendingItems); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasPendingItems {
		return nil, fmt.Errorf("CONFLICT")
	}

	var latestOrderID int
	orderQuery := `
		SELECT order_id
		FROM customer_orders
		WHERE session_id = $1
		ORDER BY order_time DESC, order_id DESC
		LIMIT 1
	`
	if err := tx.QueryRow(ctx, orderQuery, req.SessionID).Scan(&latestOrderID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("UNPROCESSABLE")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	var totalAmount float64
	totalQuery := `
		SELECT COALESCE(SUM(oi.quantity * oi.unit_price), 0)
		FROM customer_orders co
		JOIN order_items oi ON oi.order_id = co.order_id
		WHERE co.session_id = $1
		  AND oi.item_status <> 'CANCELLED'
	`
	if err := tx.QueryRow(ctx, totalQuery, req.SessionID).Scan(&totalAmount); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if totalAmount <= 0 {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}
	if req.ReceivedAmount < totalAmount {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	var paymentID int
	createPaymentQuery := `
		INSERT INTO payments (order_id, payment_method_id, total_amount, payment_status)
		VALUES ($1, $2, $3, 'PAID')
		RETURNING payment_id
	`
	if err := tx.QueryRow(ctx, createPaymentQuery, latestOrderID, req.PaymentMethodID, totalAmount).Scan(&paymentID); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	var receiptID int
	var issueDate time.Time
	if err := tx.QueryRow(ctx, `
		INSERT INTO receipts (payment_id, receipt_number, total_amount)
		VALUES ($1, $2, $3)
		RETURNING receipt_id, issue_date
	`, paymentID, fmt.Sprintf("TMP-%d", paymentID), totalAmount).Scan(&receiptID, &issueDate); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	receiptNumber := fmt.Sprintf("RCT-%s-%04d", issueDate.Format("20060102"), receiptID)
	if _, err := tx.Exec(ctx, `UPDATE receipts SET receipt_number = $1 WHERE receipt_id = $2`, receiptNumber, receiptID); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	var finalSessionStatus string
	if err := tx.QueryRow(ctx, `
		UPDATE table_sessions
		SET session_status = 'CLOSED', end_time = CURRENT_TIMESTAMP
		WHERE session_id = $1
		RETURNING session_status
	`, req.SessionID).Scan(&finalSessionStatus); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	var finalTableStatus string
	if err := tx.QueryRow(ctx, `
		UPDATE restaurant_tables
		SET table_status = 'AVAILABLE'
		WHERE table_id = $1
		RETURNING table_status
	`, tableID).Scan(&finalTableStatus); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return &models.CashierCheckoutResult{
		PaymentID:     paymentID,
		ReceiptID:     receiptID,
		ReceiptNumber: receiptNumber,
		SessionID:     req.SessionID,
		SessionStatus: finalSessionStatus,
		TableID:       tableID,
		TableStatus:   finalTableStatus,
		ChangeAmount:  req.ReceivedAmount - totalAmount,
	}, nil
}
