package services

import (
	"context"
	"fmt"
	"time"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type MenuService struct {
	repo          *repositories.MenuRepository
	qrSessionRepo *repositories.QRSessionRepository
}

func NewMenuService(repo *repositories.MenuRepository, qrSessionRepo *repositories.QRSessionRepository) *MenuService {
	return &MenuService{repo: repo, qrSessionRepo: qrSessionRepo}
}

func (s *MenuService) GetCustomerMenus(ctx context.Context, qrToken string) (*models.CustomerMenuResponse, error) {
	qr, err := s.qrSessionRepo.GetByToken(ctx, qrToken)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	if time.Now().After(qr.ExpiredAt) {
		return nil, fmt.Errorf("GONE")
	}

	if qr.SessionStatus == "CLOSED" {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	resp, err := s.repo.GetCustomerMenus(ctx, qr.TableID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
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

func (s *MenuService) UpdateStatus(ctx context.Context, menuID int, status bool) (*models.UpdateMenuStatusResponse, error) {
	resp, err := s.repo.UpdateStatus(ctx, menuID, status)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return nil, fmt.Errorf("NOT_FOUND")
		default:
			return nil, fmt.Errorf("INTERNAL")
		}
	}
	return resp, nil
}
