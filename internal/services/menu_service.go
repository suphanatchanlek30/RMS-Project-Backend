package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type MenuService struct {
	repo *repositories.MenuRepository
}

func NewMenuService(repo *repositories.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) GetCustomerMenus(ctx context.Context) ([]models.Menu, error) {
	return s.repo.GetCustomerMenus(ctx)
}

func (s *MenuService) Create(ctx context.Context, req models.CreateMenuRequest) (*models.CreateMenuResponse, error) {
	resp, err := s.repo.Create(ctx, req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return nil, fmt.Errorf("NOT_FOUND")
		case "CONFLICT":
			return nil, fmt.Errorf("CONFLICT")
		default:
			return nil, fmt.Errorf("INTERNAL")
		}
	}

	return resp, nil
}
