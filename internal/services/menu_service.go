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

func (s *MenuService) GetAll(ctx context.Context, categoryID *int, keyword string, status *bool, page, limit int) ([]models.Menu, int, error) {
	items, total, err := s.repo.GetAll(ctx, categoryID, keyword, status, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("INTERNAL")
	}
	return items, total, nil
}

func (s *MenuService) GetByID(ctx context.Context, menuID int) (*models.MenuDetail, error) {
	m, err := s.repo.GetByID(ctx, menuID)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return nil, fmt.Errorf("NOT_FOUND")
		default:
			return nil, fmt.Errorf("INTERNAL")
		}
	}
	return m, nil
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

func (s *MenuService) Update(ctx context.Context, menuID int, req models.UpdateMenuRequest) (*models.UpdateMenuResponse, error) {
	resp, err := s.repo.Update(ctx, menuID, req)
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
