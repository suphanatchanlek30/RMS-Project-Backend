package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type CashierService struct {
	repo *repositories.CashierRepository
}

func NewCashierService(repo *repositories.CashierRepository) *CashierService {
	return &CashierService{repo: repo}
}

func (s *CashierService) GetTablesOverview(ctx context.Context) ([]models.CashierTableOverviewItem, error) {
	items, err := s.repo.GetTablesOverview(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	return items, nil
}

func (s *CashierService) GetSessionCheckout(ctx context.Context, sessionID int) (*models.SessionCheckoutResponse, error) {
	if sessionID <= 0 {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	resp, err := s.repo.GetSessionCheckout(ctx, sessionID)
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

func (s *CashierService) Checkout(ctx context.Context, req models.CashierCheckoutRequest) (*models.CashierCheckoutResult, error) {
	if req.SessionID <= 0 || req.PaymentMethodID <= 0 || req.ReceivedAmount < 0 {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	resp, err := s.repo.ProcessCheckout(ctx, req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND_SESSION":
			return nil, fmt.Errorf("NOT_FOUND")
		case "NOT_FOUND_PAYMENT_METHOD":
			return nil, fmt.Errorf("NOT_FOUND")
		case "CONFLICT":
			return nil, fmt.Errorf("CONFLICT")
		case "UNPROCESSABLE":
			return nil, fmt.Errorf("UNPROCESSABLE")
		default:
			return nil, fmt.Errorf("INTERNAL")
		}
	}

	return resp, nil
}
