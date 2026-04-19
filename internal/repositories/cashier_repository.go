package repositories

import (
	"context"
	"fmt"
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

func (r *CashierRepository) Checkout(ctx context.Context, req *models.CheckoutRequest) (*models.CheckoutResponseData, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer tx.Rollback(ctx)

	// ตรวจสอบ session และ table
	var sessionID, tableID int
	var sessionStatus string
	err = tx.QueryRow(ctx, `
		SELECT ts.session_id, ts.table_id, ts.session_status
		FROM table_sessions ts
		WHERE ts.session_id = $1
	`, req.SessionID).Scan(&sessionID, &tableID, &sessionStatus)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}
	if sessionStatus != "OPEN" {
		return nil, fmt.Errorf("CONFLICT")
	}

	// ตรวจสอบว่ามี payment แล้วหรือไม่
	var hasPaid bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM payments p
			JOIN customer_orders co ON co.order_id = p.order_id
			WHERE co.session_id = $1 AND p.payment_status = 'PAID'
		)
	`, req.SessionID).Scan(&hasPaid)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasPaid {
		return nil, fmt.Errorf("CONFLICT")
	}

	// คำนวณ total amount
	var totalAmount float64
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(oi.quantity * oi.unit_price), 0)
		FROM customer_orders co
		JOIN order_items oi ON oi.order_id = co.order_id
		WHERE co.session_id = $1 AND oi.item_status <> 'CANCELLED'
	`, req.SessionID).Scan(&totalAmount)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	// ตรวจสอบ receivedAmount
	if req.ReceivedAmount < totalAmount {
		return nil, fmt.Errorf("VALIDATION")
	}

	// ได้ order_id ล่าสุด
	var orderID int
	err = tx.QueryRow(ctx, `
		SELECT co.order_id
		FROM customer_orders co
		WHERE co.session_id = $1
		ORDER BY co.order_time DESC, co.order_id DESC
		LIMIT 1
	`, req.SessionID).Scan(&orderID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	// ได้ payment method name
	var paymentMethodName string
	err = tx.QueryRow(ctx, `
		SELECT method_name
		FROM payment_methods
		WHERE payment_method_id = $1
	`, req.PaymentMethodID).Scan(&paymentMethodName)
	if err != nil {
		return nil, fmt.Errorf("VALIDATION")
	}

	// สร้าง payment
	var paymentID int
	err = tx.QueryRow(ctx, `
		INSERT INTO payments (order_id, session_id, payment_method_id, total_amount, payment_status)
		VALUES ($1, $2, $3, $4, 'PAID')
		RETURNING payment_id
	`, orderID, req.SessionID, req.PaymentMethodID, totalAmount).Scan(&paymentID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	// สร้าง receipt
	var receiptID int
	var issueDate time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO receipts (payment_id, receipt_number, total_amount)
		VALUES ($1, 'TMP', $2)
		RETURNING receipt_id, issue_date
	`, paymentID, totalAmount).Scan(&receiptID, &issueDate)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	receiptNumber := fmt.Sprintf("RCT-%s-%04d", issueDate.Format("20060102"), receiptID)
	_, err = tx.Exec(ctx, `
		UPDATE receipts SET receipt_number = $1 WHERE receipt_id = $2
	`, receiptNumber, receiptID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	// ปิด session
	_, err = tx.Exec(ctx, `
		UPDATE table_sessions
		SET session_status = 'CLOSED', end_time = CURRENT_TIMESTAMP
		WHERE session_id = $1
	`, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	// update table status
	_, err = tx.Exec(ctx, `
		UPDATE restaurant_tables
		SET table_status = 'AVAILABLE'
		WHERE table_id = $1
	`, tableID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	resp := &models.CheckoutResponseData{
		PaymentID:     paymentID,
		ReceiptID:     receiptID,
		ReceiptNumber: receiptNumber,
		SessionID:     req.SessionID,
		SessionStatus: "CLOSED",
		TableID:       tableID,
		TableStatus:   "AVAILABLE",
		ChangeAmount:  req.ReceivedAmount - totalAmount,
	}

	return resp, nil
}