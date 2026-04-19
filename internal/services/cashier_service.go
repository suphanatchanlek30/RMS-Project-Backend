package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type CashierService struct {
	repo             *repositories.CashierRepository
	tableSessionRepo *repositories.TableSessionRepository
}

func NewCashierService(repo *repositories.CashierRepository, tableSessionRepo *repositories.TableSessionRepository) *CashierService {
	return &CashierService{repo: repo, tableSessionRepo: tableSessionRepo}
}

func (s *CashierService) GetTablesOverview(ctx context.Context) ([]models.CashierTableOverview, error) {
	return s.repo.GetTablesOverview(ctx)
}

func (s *CashierService) GetCheckout(ctx context.Context, sessionID int) (*models.CheckoutResponse, error) {
	return s.repo.GetCheckout(ctx, sessionID)
}

func (s *CashierService) Checkout(ctx context.Context, req *models.CheckoutRequest) (*models.CheckoutResponseData, error) {
	if req.SessionID <= 0 || req.PaymentMethodID <= 0 || req.ReceivedAmount < 0 {
		return nil, fmt.Errorf("VALIDATION")
	}

	hasUnbillable, err := s.tableSessionRepo.HasUnbillableItems(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasUnbillable {
		return nil, fmt.Errorf("NOT_READY")
	}

	return s.repo.Checkout(ctx, req)
}