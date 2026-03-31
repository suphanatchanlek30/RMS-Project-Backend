package repositories

import (
	"context"
	"errors"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account disabled")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) FindEmployeeByEmail(ctx context.Context, email string) (*models.EmployeeAuth, error) {
	query := `
		SELECT
			e.employee_id,
			e.employee_name,
			e.email,
			e.password_hash,
			e.employee_status,
			e.role_id,
			r.role_name
		FROM employees e
		JOIN roles r ON e.role_id = r.role_id
		WHERE e.email = $1
		LIMIT 1
	`

	var employee models.EmployeeAuth
	if err := r.DB.QueryRow(ctx, query, email).Scan(
		&employee.EmployeeID,
		&employee.EmployeeName,
		&employee.Email,
		&employee.PasswordHash,
		&employee.EmployeeStatus,
		&employee.RoleID,
		&employee.RoleName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &employee, nil
}

func (r *AuthRepository) FindEmployeeByID(ctx context.Context, employeeID int) (*models.EmployeeAuth, error) {
	query := `
		SELECT
			e.employee_id,
			e.employee_name,
			e.email,
			COALESCE(e.phone_number, '') AS phone_number,
			e.employee_status,
			e.role_id,
			r.role_name
		FROM employees e
		JOIN roles r ON e.role_id = r.role_id
		WHERE e.employee_id = $1
		LIMIT 1
	`

	var employee models.EmployeeAuth
	if err := r.DB.QueryRow(ctx, query, employeeID).Scan(
		&employee.EmployeeID,
		&employee.EmployeeName,
		&employee.Email,
		&employee.PhoneNumber,
		&employee.EmployeeStatus,
		&employee.RoleID,
		&employee.RoleName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	return &employee, nil
}
