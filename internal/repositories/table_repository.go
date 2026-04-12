package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TableRepository struct {
	DB *pgxpool.Pool
}

func NewTableRepository(db *pgxpool.Pool) *TableRepository {
	return &TableRepository{DB: db}
}

func (r *TableRepository) GetAll(ctx context.Context, status string, page, limit int) ([]models.RestaurantTable, error) {
	offset := (page - 1) * limit

	query := `
		SELECT table_id, table_number, capacity, table_status, created_at
		FROM restaurant_tables
		WHERE ($1 = '' OR table_status = $1)
		ORDER BY table_number ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(ctx, query, status, limit, offset)
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

	return tables, rows.Err()
}

func (r *TableRepository) GetByID(ctx context.Context, tableID int) (*models.RestaurantTable, error) {
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

func (r *TableRepository) Create(ctx context.Context, tableNumber string, capacity int) (*models.RestaurantTable, error) {
	query := `
		INSERT INTO restaurant_tables (table_number, capacity, table_status, created_at)
		VALUES ($1, $2, 'AVAILABLE', CURRENT_TIMESTAMP)
		RETURNING table_id, table_number, capacity, table_status, created_at
	`

	var t models.RestaurantTable
	err := r.DB.QueryRow(ctx, query, tableNumber, capacity).Scan(
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

func (r *TableRepository) Update(ctx context.Context, tableId int, req models.UpdateTableRequest) (*models.RestaurantTable, error) {
	query := `
		UPDATE restaurant_tables
		SET table_number = $1,
		    capacity = $2
		WHERE table_id = $3
		RETURNING table_id, table_number, capacity, table_status, created_at
	`

	var table models.RestaurantTable
	err := r.DB.QueryRow(ctx, query,
		req.TableNumber,
		req.Capacity,
		tableId,
	).Scan(
		&table.TableID,
		&table.TableNumber,
		&table.Capacity,
		&table.TableStatus,
		&table.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("NOT_FOUND")
		}

		if strings.Contains(err.Error(), "duplicate") {
			return nil, fmt.Errorf("DUPLICATE")
		}

		return nil, err
	}

	return &table, nil
}
