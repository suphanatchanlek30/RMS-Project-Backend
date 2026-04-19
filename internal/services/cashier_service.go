package services

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type CashierService struct {
	repo *repositories.CashierRepository
}

func NewCashierService(repo *repositories.CashierRepository) *CashierService {
	return &CashierService{repo: repo}
}

func (s *CashierService) GetTablesOverview(ctx context.Context) ([]models.CashierTableOverview, error) {
	return s.repo.GetTablesOverview(ctx)
}

func (s *CashierService) GetCheckout(ctx context.Context, sessionID int) (*models.CheckoutResponse, error) {
	return s.repo.GetCheckout(ctx, sessionID)
}

func (s *CashierService) Checkout(ctx context.Context, req *models.CheckoutRequest) (*models.CheckoutResponseData, error) {
	return s.repo.Checkout(ctx, req)
}