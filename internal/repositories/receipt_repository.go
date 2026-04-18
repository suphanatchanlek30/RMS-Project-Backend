package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type ReceiptRepository struct {
	DB *pgxpool.Pool
}

func NewReceiptRepository(db *pgxpool.Pool) *ReceiptRepository {
	return &ReceiptRepository{DB: db}
}

func (r *ReceiptRepository) CreateForPayment(ctx context.Context, paymentID int, totalAmount float64) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("INTERNAL")
	}
	defer tx.Rollback(ctx)

	var existingID int
	err = tx.QueryRow(ctx, `SELECT receipt_id FROM receipts WHERE payment_id = $1`, paymentID).Scan(&existingID)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("INTERNAL")
		}
		return nil
	}
	if err != pgx.ErrNoRows {
		return fmt.Errorf("INTERNAL")
	}

	tmpNumber := fmt.Sprintf("TMP-%d", paymentID)
	var receiptID int
	var issueDate time.Time
	err = tx.QueryRow(ctx,
		`INSERT INTO receipts (payment_id, receipt_number, total_amount)
		 VALUES ($1, $2, $3)
		 RETURNING receipt_id, issue_date`,
		paymentID, tmpNumber, totalAmount,
	).Scan(&receiptID, &issueDate)
	if err != nil {
		return fmt.Errorf("INTERNAL")
	}

	receiptNumber := fmt.Sprintf("RCT-%s-%04d", issueDate.Format("20060102"), receiptID)
	if _, err := tx.Exec(ctx, `UPDATE receipts SET receipt_number = $1 WHERE receipt_id = $2`, receiptNumber, receiptID); err != nil {
		return fmt.Errorf("INTERNAL")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("INTERNAL")
	}

	return nil
}

func (r *ReceiptRepository) GetByPaymentID(ctx context.Context, paymentID int) (*models.ReceiptDetailResponse, error) {
	query := `
		SELECT
			r.receipt_id,
			r.receipt_number,
			r.issue_date,
			r.total_amount,
			p.payment_id,
			pm.method_name,
			p.payment_time,
			t.table_id,
			t.table_number,
			co.order_id
		FROM receipts r
		JOIN payments p ON p.payment_id = r.payment_id
		JOIN payment_methods pm ON pm.payment_method_id = p.payment_method_id
		JOIN customer_orders co ON co.order_id = p.order_id
		JOIN restaurant_tables t ON t.table_id = co.table_id
		WHERE p.payment_id = $1
	`

	resp := &models.ReceiptDetailResponse{}
	var orderID int
	if err := r.DB.QueryRow(ctx, query, paymentID).Scan(
		&resp.ReceiptID,
		&resp.ReceiptNumber,
		&resp.IssueDate,
		&resp.TotalAmount,
		&resp.Payment.PaymentID,
		&resp.Payment.PaymentMethodName,
		&resp.Payment.PaymentTime,
		&resp.Table.TableID,
		&resp.Table.TableNumber,
		&orderID,
	); err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	items, err := r.getReceiptItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	resp.Items = items

	return resp, nil
}

func (r *ReceiptRepository) GetByReceiptID(ctx context.Context, receiptID int) (*models.ReceiptDetailResponse, error) {
	query := `
		SELECT
			r.receipt_id,
			r.receipt_number,
			r.issue_date,
			r.total_amount,
			p.payment_id,
			pm.method_name,
			p.payment_time,
			t.table_id,
			t.table_number,
			co.order_id
		FROM receipts r
		JOIN payments p ON p.payment_id = r.payment_id
		JOIN payment_methods pm ON pm.payment_method_id = p.payment_method_id
		JOIN customer_orders co ON co.order_id = p.order_id
		JOIN restaurant_tables t ON t.table_id = co.table_id
		WHERE r.receipt_id = $1
	`

	resp := &models.ReceiptDetailResponse{}
	var orderID int
	if err := r.DB.QueryRow(ctx, query, receiptID).Scan(
		&resp.ReceiptID,
		&resp.ReceiptNumber,
		&resp.IssueDate,
		&resp.TotalAmount,
		&resp.Payment.PaymentID,
		&resp.Payment.PaymentMethodName,
		&resp.Payment.PaymentTime,
		&resp.Table.TableID,
		&resp.Table.TableNumber,
		&orderID,
	); err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	items, err := r.getReceiptItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	resp.Items = items

	return resp, nil
}

func (r *ReceiptRepository) getReceiptItems(ctx context.Context, orderID int) ([]models.ReceiptItem, error) {
	query := `
		SELECT
			m.menu_name,
			oi.quantity,
			oi.unit_price,
			(oi.quantity * oi.unit_price) AS line_total
		FROM order_items oi
		JOIN menus m ON m.menu_id = oi.menu_id
		WHERE oi.order_id = $1
		  AND oi.item_status <> 'CANCELLED'
		ORDER BY oi.order_item_id ASC
	`

	rows, err := r.DB.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.ReceiptItem, 0)
	for rows.Next() {
		var item models.ReceiptItem
		if err := rows.Scan(&item.MenuName, &item.Quantity, &item.UnitPrice, &item.LineTotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
