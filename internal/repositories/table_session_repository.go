package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type TableSessionRepository struct {
	DB *pgxpool.Pool
}

func NewTableSessionRepository(db *pgxpool.Pool) *TableSessionRepository {
	return &TableSessionRepository{DB: db}
}

func (r *TableSessionRepository) GetTableByID(ctx context.Context, tableID int) (*models.RestaurantTable, error) {
	query := `
		SELECT table_id, table_number, capacity, table_status, created_at
		FROM restaurant_tables
		WHERE table_id = $1
	`

	var t models.RestaurantTable
	err := r.DB.QueryRow(ctx, query, tableID).Scan(
		&t.TableID,
		&t.TableNumber,
		&t.Capacity,
		&t.TableStatus,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TableSessionRepository) HasOpenSession(ctx context.Context, tableID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM table_sessions
			WHERE table_id = $1 AND session_status = 'OPEN'
		)
	`

	var exists bool
	err := r.DB.QueryRow(ctx, query, tableID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TableSessionRepository) OpenSession(ctx context.Context, tableID int) (*models.OpenTableResponse, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer tx.Rollback(ctx)

	insertQuery := `
		INSERT INTO table_sessions (table_id, start_time, session_status)
		VALUES ($1, CURRENT_TIMESTAMP, 'OPEN')
		RETURNING session_id, table_id, start_time, session_status
	`

	var resp models.OpenTableResponse
	err = tx.QueryRow(ctx, insertQuery, tableID).Scan(
		&resp.SessionID,
		&resp.TableID,
		&resp.StartTime,
		&resp.SessionStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	updateQuery := `
		UPDATE restaurant_tables
		SET table_status = 'OCCUPIED'
		WHERE table_id = $1
		RETURNING table_number
	`

	err = tx.QueryRow(ctx, updateQuery, tableID).Scan(&resp.TableNumber)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return &resp, nil
}

func (r *TableSessionRepository) GetByID(ctx context.Context, sessionID int) (*models.TableSessionDetail, error) {
	query := `
		SELECT ts.session_id, ts.table_id, rt.table_number, ts.start_time, ts.end_time, ts.session_status
		FROM table_sessions ts
		JOIN restaurant_tables rt ON ts.table_id = rt.table_id
		WHERE ts.session_id = $1
	`

	var s models.TableSessionDetail
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(
		&s.SessionID,
		&s.TableID,
		&s.TableNumber,
		&s.StartTime,
		&s.EndTime,
		&s.SessionStatus,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *TableSessionRepository) GetCurrentSessionByTableID(ctx context.Context, tableID int) (*models.CurrentSessionResponse, error) {
	query := `
		SELECT session_id, table_id, session_status, start_time
		FROM table_sessions
		WHERE table_id = $1 AND session_status = 'OPEN'
		ORDER BY start_time DESC
		LIMIT 1
	`

	var s models.CurrentSessionResponse
	err := r.DB.QueryRow(ctx, query, tableID).Scan(
		&s.SessionID,
		&s.TableID,
		&s.SessionStatus,
		&s.StartTime,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *TableSessionRepository) GetSessionByID(ctx context.Context, sessionID int) (*models.TableSession, error) {
	query := `
		SELECT session_id, table_id, start_time, end_time, session_details, session_status
		FROM table_sessions
		WHERE session_id = $1
	`

	var s models.TableSession
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(
		&s.SessionID,
		&s.TableID,
		&s.StartTime,
		&s.EndTime,
		&s.SessionDetails,
		&s.SessionStatus,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *TableSessionRepository) HasPendingOrders(ctx context.Context, sessionID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM customer_orders
			WHERE session_id = $1 AND order_status IN ('PENDING', 'PREPARING')
		)
	`

	var exists bool
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TableSessionRepository) CloseSession(ctx context.Context, sessionID int, tableID int) (*models.CloseSessionResponse, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer tx.Rollback(ctx)

	updateSessionQuery := `
		UPDATE table_sessions
		SET session_status = 'CLOSED', end_time = CURRENT_TIMESTAMP
		WHERE session_id = $1
		RETURNING session_id, session_status, end_time
	`

	var resp models.CloseSessionResponse
	err = tx.QueryRow(ctx, updateSessionQuery, sessionID).Scan(
		&resp.SessionID,
		&resp.SessionStatus,
		&resp.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	updateTableQuery := `
		UPDATE restaurant_tables
		SET table_status = 'AVAILABLE'
		WHERE table_id = $1
		RETURNING table_id, table_status
	`

	err = tx.QueryRow(ctx, updateTableQuery, tableID).Scan(
		&resp.TableID,
		&resp.TableStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return &resp, nil
}

func (r *TableSessionRepository) HasUnbillableItems(ctx context.Context, sessionID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM customer_orders co
			JOIN order_items oi ON oi.order_id = co.order_id
			WHERE co.session_id = $1
			  AND oi.item_status IN ('WAITING', 'PREPARING')
		)
	`

	var exists bool
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TableSessionRepository) GetSessionBill(ctx context.Context, sessionID int) (*models.SessionBillResponse, error) {
	infoQuery := `
		SELECT ts.session_id, ts.table_id, rt.table_number
		FROM table_sessions ts
		JOIN restaurant_tables rt ON rt.table_id = ts.table_id
		WHERE ts.session_id = $1
	`

	resp := &models.SessionBillResponse{}
	err := r.DB.QueryRow(ctx, infoQuery, sessionID).Scan(
		&resp.SessionID,
		&resp.TableID,
		&resp.TableNumber,
	)
	if err != nil {
		return nil, err
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

	rows, err := r.DB.Query(ctx, itemsQuery, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.BillItemResponse, 0)
	var subtotal float64

	for rows.Next() {
		var item models.BillItemResponse
		if err := rows.Scan(
			&item.OrderItemID,
			&item.MenuName,
			&item.Quantity,
			&item.UnitPrice,
			&item.LineTotal,
		); err != nil {
			return nil, err
		}

		subtotal += item.LineTotal
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("NO_ITEMS")
	}

	resp.Items = items
	resp.Subtotal = subtotal
	resp.ServiceCharge = 0
	resp.VAT = 0
	resp.TotalAmount = subtotal + resp.ServiceCharge + resp.VAT

	return resp, nil
}
