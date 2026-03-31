package services

import (
	"context"
	"errors"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
	"github.com/suphanatchanlek30/rms-project-backend/internal/utils"
)

type AuthService struct {
	repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponseData, error) {
	employee, err := s.repo.FindEmployeeByEmail(ctx, req.Email)
	if err != nil {
		return nil, repositories.ErrInvalidCredentials
	}

	if !employee.EmployeeStatus {
		return nil, repositories.ErrAccountDisabled
	}

	if err := utils.CheckPasswordHash(req.Password, employee.PasswordHash); err != nil {
		return nil, repositories.ErrInvalidCredentials
	}

	token, expiresIn, err := utils.GenerateJWT(
		employee.EmployeeID,
		employee.RoleID,
		employee.RoleName,
		employee.Email,
	)
	if err != nil {
		return nil, err
	}

	resp := &models.LoginResponseData{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Employee: models.LoginEmployeeResponse{
			EmployeeID:   employee.EmployeeID,
			EmployeeName: employee.EmployeeName,
			Email:        employee.Email,
			Role: models.RoleResponse{
				RoleID:   employee.RoleID,
				RoleName: employee.RoleName,
			},
		},
	}

	return resp, nil
}

func (s *AuthService) GetMe(ctx context.Context, employeeID int) (*models.MeResponseData, error) {
	employee, err := s.repo.FindEmployeeByID(ctx, employeeID)
	if err != nil {
		return nil, repositories.ErrInvalidToken
	}

	if !employee.EmployeeStatus {
		return nil, repositories.ErrAccountDisabled
	}

	resp := &models.MeResponseData{
		EmployeeID:     employee.EmployeeID,
		EmployeeName:   employee.EmployeeName,
		Email:          employee.Email,
		PhoneNumber:    employee.PhoneNumber,
		EmployeeStatus: employee.EmployeeStatus,
		Role: models.RoleResponse{
			RoleID:   employee.RoleID,
			RoleName: employee.RoleName,
		},
	}

	return resp, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	return nil
}

var ErrUnauthorized = errors.New("unauthorized")
