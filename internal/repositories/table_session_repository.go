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
