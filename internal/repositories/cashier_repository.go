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