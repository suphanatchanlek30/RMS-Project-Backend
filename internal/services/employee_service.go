package services

import (
	"context"
	"errors"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeService(r *repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: r}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, req models.CreateEmployeeRequest) (*models.Employee, error) {

	exists, err := s.repo.CheckEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("EMAIL_EXISTS")
	}

	roleExists, err := s.repo.CheckRole(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, errors.New("ROLE_NOT_FOUND")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateEmployee(ctx, req, string(hash))
}

func (s *EmployeeService) GetEmployees(
	ctx context.Context,
	roleID *int,
	status *bool,
	search string,
	page int,
	limit int,
) ([]models.Employee, int, error) {

	return s.repo.GetEmployees(ctx, roleID, status, search, page, limit)
}

func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id int) (*models.Employee, error) {
	emp, err := s.repo.GetEmployeeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if emp == nil {
		return nil, errors.New("NOT_FOUND")
	}

	return emp, nil
}

func (s *EmployeeService) UpdateEmployee(
	ctx context.Context,
	id int,
	req models.UpdateEmployeeRequest,
) (*models.Employee, error) {

	roleExists, err := s.repo.CheckRole(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, errors.New("ROLE_NOT_FOUND")
	}

	phoneExists, err := s.repo.CheckPhoneDuplicate(ctx, id, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, errors.New("PHONE_EXISTS")
	}

	emp, err := s.repo.UpdateEmployee(ctx, id, req)
	if err != nil {
		return nil, err
	}

	if emp == nil {
		return nil, errors.New("NOT_FOUND")
	}

	return emp, nil
}
