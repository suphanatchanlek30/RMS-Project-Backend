package repositories

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuRepository struct {
	DB *pgxpool.Pool
}

func NewMenuRepository(db *pgxpool.Pool) *MenuRepository {
	return &MenuRepository{DB: db}
}

func (r *MenuRepository) GetCustomerMenus(ctx context.Context) ([]models.Menu, error) {
	query := `
		SELECT
			m.menu_id,
			m.menu_name,
			m.category_id,
			c.category_name,
			m.price,
			COALESCE(m.description, '') AS description,
			m.menu_status,
			m.created_at
		FROM menus m
		JOIN menu_categories c ON m.category_id = c.category_id
		WHERE m.menu_status = TRUE
		ORDER BY c.category_name ASC, m.menu_name ASC
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var m models.Menu
		if err := rows.Scan(
			&m.MenuID,
			&m.MenuName,
			&m.CategoryID,
			&m.CategoryName,
			&m.Price,
			&m.Description,
			&m.MenuStatus,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}
