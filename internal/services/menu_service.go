package services

import (
	"context"

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
