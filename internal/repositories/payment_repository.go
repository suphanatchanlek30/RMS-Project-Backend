package repositories

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type PaymentRepository struct {
	DB *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) GetPaymentMethodNameByID(ctx context.Context, paymentMethodID int) (string, error) {
	query := `
		SELECT method_name
		FROM payment_methods
		WHERE payment_method_id = $1
	`

	var methodName string
	if err := r.DB.QueryRow(ctx, query, paymentMethodID).Scan(&methodName); err != nil {
		return "", fmt.Errorf("NOT_FOUND")
	}

	return methodName, nil
}

func (r *PaymentRepository) GetLatestOrderIDBySession(ctx context.Context, sessionID int) (int, error) {
	query := `
		SELECT co.order_id
		FROM customer_orders co
		WHERE co.session_id = $1
		ORDER BY co.order_time DESC, co.order_id DESC
		LIMIT 1
	`

	var orderID int
	if err := r.DB.QueryRow(ctx, query, sessionID).Scan(&orderID); err != nil {
		return 0, fmt.Errorf("NOT_FOUND")
	}

	return orderID, nil
}

func (r *PaymentRepository) HasPaidPaymentBySession(ctx context.Context, sessionID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM payments p
			JOIN customer_orders co ON co.order_id = p.order_id
			WHERE co.session_id = $1
			  AND p.payment_status = 'PAID'
		)
	`

	var exists bool
	if err := r.DB.QueryRow(ctx, query, sessionID).Scan(&exists); err != nil {
		return false, fmt.Errorf("INTERNAL")
	}

	return exists, nil
}

func (r *PaymentRepository) CreatePayment(
	ctx context.Context,
	orderID int,
	sessionID int,
	paymentMethodID int,
	paymentMethodName string,
	totalAmount float64,
	receivedAmount float64,
) (*models.CreatePaymentResponse, error) {
	query := `
		INSERT INTO payments (order_id, payment_method_id, total_amount, payment_status)
		VALUES ($1, $2, $3, 'PAID')
		RETURNING payment_id, payment_time, payment_status
	`

	resp := &models.CreatePaymentResponse{
		SessionID:         sessionID,
		PaymentMethodID:   paymentMethodID,
		PaymentMethodName: paymentMethodName,
		TotalAmount:       totalAmount,
		ReceivedAmount:    receivedAmount,
		ChangeAmount:      receivedAmount - totalAmount,
	}

	if err := r.DB.QueryRow(ctx, query, orderID, paymentMethodID, totalAmount).Scan(
		&resp.PaymentID,
		&resp.PaymentTime,
		&resp.PaymentStatus,
	); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, paymentID int) (*models.PaymentDetailResponse, error) {
	query := `
		SELECT
			p.payment_id,
			co.session_id,
			p.payment_method_id,
			pm.method_name,
			p.total_amount,
			p.payment_time,
			p.payment_status
		FROM payments p
		JOIN customer_orders co ON co.order_id = p.order_id
		JOIN payment_methods pm ON pm.payment_method_id = p.payment_method_id
		WHERE p.payment_id = $1
	`

	var resp models.PaymentDetailResponse
	if err := r.DB.QueryRow(ctx, query, paymentID).Scan(
		&resp.PaymentID,
		&resp.SessionID,
		&resp.PaymentMethodID,
		&resp.PaymentMethodName,
		&resp.TotalAmount,
		&resp.PaymentTime,
		&resp.PaymentStatus,
	); err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	return &resp, nil
}

func (r *PaymentRepository) GetAll(ctx context.Context, filter models.PaymentListFilter) ([]models.PaymentListItem, int, error) {
	query := `
		SELECT
			p.payment_id,
			co.session_id,
			pm.method_name,
			p.total_amount,
			p.payment_status,
			p.payment_time
		FROM payments p
		JOIN customer_orders co ON co.order_id = p.order_id
		JOIN payment_methods pm ON pm.payment_method_id = p.payment_method_id
		WHERE 1=1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM payments p
		JOIN customer_orders co ON co.order_id = p.order_id
		JOIN payment_methods pm ON pm.payment_method_id = p.payment_method_id
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if filter.DateFrom != "" {
		clause := " AND p.payment_time >= $" + strconv.Itoa(argPos)
		query += clause
		countQuery += clause
		args = append(args, filter.DateFrom)
		argPos++
	}

	if filter.DateTo != "" {
		clause := " AND p.payment_time <= $" + strconv.Itoa(argPos)
		query += clause
		countQuery += clause
		args = append(args, filter.DateTo)
		argPos++
	}

	if filter.PaymentMethodID != nil {
		clause := " AND p.payment_method_id = $" + strconv.Itoa(argPos)
		query += clause
		countQuery += clause
		args = append(args, *filter.PaymentMethodID)
		argPos++
	}

	if filter.Status != "" {
		clause := " AND p.payment_status = $" + strconv.Itoa(argPos)
		query += clause
		countQuery += clause
		args = append(args, filter.Status)
		argPos++
	}

	var total int
	if err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("INTERNAL")
	}

	offset := (filter.Page - 1) * filter.Limit
	query += " ORDER BY p.payment_time DESC"
	query += " LIMIT $" + strconv.Itoa(argPos)
	args = append(args, filter.Limit)
	argPos++
	query += " OFFSET $" + strconv.Itoa(argPos)
	args = append(args, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	items := make([]models.PaymentListItem, 0)
	for rows.Next() {
		var item models.PaymentListItem
		if err := rows.Scan(
			&item.PaymentID,
			&item.SessionID,
			&item.PaymentMethodName,
			&item.TotalAmount,
			&item.PaymentStatus,
			&item.PaymentTime,
		); err != nil {
			return nil, 0, fmt.Errorf("INTERNAL")
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("INTERNAL")
	}

	return items, total, nil
}
