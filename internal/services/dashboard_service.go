package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type DashboardService struct {
	repo *repositories.DashboardRepository
}

func NewDashboardService(repo *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetSummary(ctx context.Context) (*models.DashboardSummaryResponse, error) {
	resp, err := s.repo.GetSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	return resp, nil
}
