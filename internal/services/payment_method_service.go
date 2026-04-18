package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type PaymentMethodService struct {
	repo *repositories.PaymentMethodRepository
}

func NewPaymentMethodService(repo *repositories.PaymentMethodRepository) *PaymentMethodService {
	return &PaymentMethodService{repo: repo}
}

func (s *PaymentMethodService) GetAll(ctx context.Context) ([]models.PaymentMethodItem, error) {
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	return items, nil
}
