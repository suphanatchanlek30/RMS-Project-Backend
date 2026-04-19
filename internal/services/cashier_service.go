package services

import (
	"context"
	"rms-project-backend/internal/models"
	"rms-project-backend/internal/repositories"
)

type CashierService interface {
	GetTablesOverview(ctx context.Context) ([]models.CashierTableOverviewItem, error)
}

type cashierService struct {
	cashierRepo repositories.CashierRepository
}

func NewCashierService(cashierRepo repositories.CashierRepository) CashierService {
	return &cashierService{
		cashierRepo: cashierRepo,
	}
}

func (s *cashierService) GetTablesOverview(ctx context.Context) ([]models.CashierTableOverviewItem, error) {
	return s.cashierRepo.GetTablesOverview(ctx)
}