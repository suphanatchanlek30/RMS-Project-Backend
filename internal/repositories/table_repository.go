package repositories

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TableRepository struct {
	DB *pgxpool.Pool
}

func NewTableRepository(db *pgxpool.Pool) *TableRepository {
	return &TableRepository{DB: db}
}

func (r *TableRepository) GetAll(ctx context.Context) ([]models.RestaurantTable, error) {
	query := `
		SELECT table_id, table_number, capacity, table_status, created_at
		FROM restaurant_tables
		ORDER BY table_number ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.RestaurantTable
	for rows.Next() {
		var t models.RestaurantTable
		if err := rows.Scan(
			&t.TableID,
			&t.TableNumber,
			&t.Capacity,
			&t.TableStatus,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
