package services

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetSalesReport(ctx context.Context, dateFrom string, dateTo string, groupBy string) ([]models.SalesReportItem, error) {
	return s.repo.GetSalesReport(ctx, dateFrom, dateTo, groupBy)
}

func (s *ReportService) GetTopMenusReport(ctx context.Context, dateFrom string, dateTo string, limit int) ([]models.TopMenuReportItem, error) {
	return s.repo.GetTopMenusReport(ctx, dateFrom, dateTo, limit)
}
