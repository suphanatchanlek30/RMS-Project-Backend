package repositories

import (
	"context"

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
