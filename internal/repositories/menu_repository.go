package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

func (r *MenuRepository) GetAll(ctx context.Context, categoryID *int, keyword string, status *bool, page, limit int) ([]models.Menu, int, error) {
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
		WHERE 1=1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM menus m
		JOIN menu_categories c ON m.category_id = c.category_id
		WHERE 1=1
	`

	args := []interface{}{}
	i := 1

	if categoryID != nil {
		filter := " AND m.category_id = $" + strconv.Itoa(i)
		query += filter
		countQuery += filter
		args = append(args, *categoryID)
		i++
	}

	if keyword != "" {
		filter := " AND m.menu_name ILIKE $" + strconv.Itoa(i)
		query += filter
		countQuery += filter
		args = append(args, "%"+keyword+"%")
		i++
	}

	if status != nil {
		filter := " AND m.menu_status = $" + strconv.Itoa(i)
		query += filter
		countQuery += filter
		args = append(args, *status)
		i++
	}

	// count total
	var total int
	err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query += " ORDER BY m.menu_id ASC"
	query += " LIMIT $" + strconv.Itoa(i)
	args = append(args, limit)
	i++
	query += " OFFSET $" + strconv.Itoa(i)
	args = append(args, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		menus = append(menus, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

func (r *MenuRepository) GetByID(ctx context.Context, menuID int) (*models.MenuDetail, error) {
	query := `
		SELECT menu_id, menu_name, category_id, price, description, menu_status
		FROM menus
		WHERE menu_id = $1
	`

	var m models.MenuDetail
	err := r.DB.QueryRow(ctx, query, menuID).Scan(
		&m.MenuID,
		&m.MenuName,
		&m.CategoryID,
		&m.Price,
		&m.Description,
		&m.MenuStatus,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return &m, nil
}

func (r *MenuRepository) Create(ctx context.Context, req models.CreateMenuRequest) (*models.CreateMenuResponse, error) {
	query := `
		INSERT INTO menus (menu_name, category_id, price, description, menu_status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING menu_id, menu_name, category_id, price, description, menu_status, created_at
	`

	var resp models.CreateMenuResponse
	err := r.DB.QueryRow(ctx, query, req.MenuName, req.CategoryID, req.Price, req.Description, req.MenuStatus).Scan(
		&resp.MenuID,
		&resp.MenuName,
		&resp.CategoryID,
		&resp.Price,
		&resp.Description,
		&resp.MenuStatus,
		&resp.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key") || strings.Contains(err.Error(), "violates foreign key") {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("CONFLICT")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return &resp, nil
}

func (r *MenuRepository) Update(ctx context.Context, menuID int, req models.UpdateMenuRequest) (*models.UpdateMenuResponse, error) {
	query := `
		UPDATE menus
		SET menu_name = $1, price = $2, description = $3
		WHERE menu_id = $4
		RETURNING menu_id, menu_name, price, description
	`

	var resp models.UpdateMenuResponse
	err := r.DB.QueryRow(ctx, query, req.MenuName, req.Price, req.Description, menuID).Scan(
		&resp.MenuID,
		&resp.MenuName,
		&resp.Price,
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
