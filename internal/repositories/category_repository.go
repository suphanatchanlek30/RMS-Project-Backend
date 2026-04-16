package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) Create(ctx context.Context, req models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	query := `
		INSERT INTO menu_categories (category_name, description)
		VALUES ($1, $2)
		RETURNING category_id, category_name, description, created_at
	`

	var resp models.CategoryResponse
	err := r.DB.QueryRow(ctx, query, req.CategoryName, req.Description).Scan(
		&resp.CategoryID,
		&resp.CategoryName,
		&resp.Description,
		&resp.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("CONFLICT")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return &resp, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]models.CategoryListItem, error) {
	query := `
		SELECT category_id, category_name, description
		FROM menu_categories
		ORDER BY category_id ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	var categories []models.CategoryListItem
	for rows.Next() {
		var c models.CategoryListItem
		if err := rows.Scan(&c.CategoryID, &c.CategoryName, &c.Description); err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}
		categories = append(categories, c)
	}

	if categories == nil {
		categories = []models.CategoryListItem{}
	}

	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, categoryID int, req models.UpdateCategoryRequest) (*models.CategoryListItem, error) {
	query := `
		UPDATE menu_categories
		SET category_name = $1, description = $2
		WHERE category_id = $3
		RETURNING category_id, category_name, description
	`

	var resp models.CategoryListItem
	err := r.DB.QueryRow(ctx, query, req.CategoryName, req.Description, categoryID).Scan(
		&resp.CategoryID,
		&resp.CategoryName,
		&resp.Description,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("CONFLICT")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return &resp, nil
}
