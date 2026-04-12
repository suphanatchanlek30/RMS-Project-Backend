package repositories

import (
	"context"
	"strconv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type EmployeeRepository struct {
	DB *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) CheckEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM employees WHERE email=$1)", email,
	).Scan(&exists)
	return exists, err
}

func (r *EmployeeRepository) CheckRole(ctx context.Context, roleID int) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM roles WHERE role_id=$1)", roleID,
	).Scan(&exists)
	return exists, err
}

func (r *EmployeeRepository) CreateEmployee(
	ctx context.Context,
	req models.CreateEmployeeRequest,
	passwordHash string,
) (*models.Employee, error) {

	var id int

	err := r.DB.QueryRow(ctx, `
		INSERT INTO employees 
		(employee_name, role_id, phone_number, email, hire_date, password_hash, employee_status)
		VALUES ($1,$2,$3,$4,$5,$6,TRUE)
		RETURNING employee_id
	`,
		req.EmployeeName,
		req.RoleID,
		req.PhoneNumber,
		req.Email,
		req.HireDate,
		passwordHash,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &models.Employee{
		EmployeeID:     id,
		EmployeeName:   req.EmployeeName,
		RoleID:         req.RoleID,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		HireDate:       req.HireDate,
		EmployeeStatus: true,
	}, nil
}

func (r *EmployeeRepository) GetEmployees(
	ctx context.Context,
	roleID *int,
	status *bool,
	search string,
	page int,
	limit int,
) ([]models.Employee, int, error) {

	query := `
	SELECT 
		e.employee_id,
		e.employee_name,
		e.role_id,
		r.role_name,
		e.phone_number,e.email,
		e.employee_status
	FROM employees e
	JOIN roles r ON e.role_id = r.role_id
	WHERE 1=1
	`

	args := []interface{}{}
	i := 1

	if roleID != nil {
		query += " AND e.role_id = $" + strconv.Itoa(i)
		args = append(args, *roleID)
		i++
	}

	if status != nil {
		query += " AND e.employee_status = $" + strconv.Itoa(i)
		args = append(args, *status)
		i++
	}

	if search != "" {
		query += " AND (e.employee_name ILIKE $" + strconv.Itoa(i) + " OR e.email ILIKE $" + strconv.Itoa(i) + ")"
		args = append(args, "%"+search+"%")
		i++
	}

	offset := (page - 1) * limit

	query += " ORDER BY e.employee_id ASC"
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

	var items []models.Employee

	for rows.Next() {
		var e models.Employee

		err := rows.Scan(
			&e.EmployeeID,
			&e.EmployeeName,
			&e.RoleID,
			&e.RoleName,
			&e.PhoneNumber,
			&e.Email,
			&e.EmployeeStatus,
		)
		if err != nil {
			return nil, 0, err
		}

		items = append(items, e)
	}

	// total
	var total int
	err = r.DB.QueryRow(ctx, "SELECT COUNT(*) FROM employees").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}