package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func ParseSalesReportQuery(dateFromRaw, dateToRaw, groupByRaw string) (models.SalesReportQuery, error) {
	dateFrom, dateTo, err := parseDateRange(dateFromRaw, dateToRaw)
	if err != nil {
		return models.SalesReportQuery{}, err
	}

	groupBy := strings.ToLower(strings.TrimSpace(groupByRaw))
	if groupBy == "" {
		groupBy = "day"
	}
	if groupBy != "day" && groupBy != "month" {
		return models.SalesReportQuery{}, fmt.Errorf("BAD_REQUEST")
	}

	return models.SalesReportQuery{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		GroupBy:  groupBy,
	}, nil
}

func ParseTopMenusReportQuery(dateFromRaw, dateToRaw string, limit int) (models.TopMenusReportQuery, error) {
	dateFrom, dateTo, err := parseDateRange(dateFromRaw, dateToRaw)
	if err != nil {
		return models.TopMenusReportQuery{}, err
	}

	if limit == 0 {
		limit = 10
	}
	if limit < 0 || limit > 100 {
		return models.TopMenusReportQuery{}, fmt.Errorf("BAD_REQUEST")
	}

	return models.TopMenusReportQuery{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
	}, nil
}

func (s *ReportService) GetSalesReport(ctx context.Context, query models.SalesReportQuery) ([]models.SalesReportItem, error) {
	items, err := s.repo.GetSalesReport(ctx, query)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return nil, fmt.Errorf("BAD_REQUEST")
		default:
			return nil, fmt.Errorf("INTERNAL")
		}
	}
	return items, nil
}

func (s *ReportService) GetTopMenusReport(ctx context.Context, query models.TopMenusReportQuery) ([]models.TopMenuReportItem, error) {
	items, err := s.repo.GetTopMenusReport(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	return items, nil
}

func parseDateRange(dateFromRaw, dateToRaw string) (time.Time, time.Time, error) {
	dateFromRaw = strings.TrimSpace(dateFromRaw)
	dateToRaw = strings.TrimSpace(dateToRaw)
	if dateFromRaw == "" || dateToRaw == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("BAD_REQUEST")
	}

	dateFrom, err := time.Parse("2006-01-02", dateFromRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("BAD_REQUEST")
	}
	dateTo, err := time.Parse("2006-01-02", dateToRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("BAD_REQUEST")
	}
	if dateTo.Before(dateFrom) {
		return time.Time{}, time.Time{}, fmt.Errorf("BAD_REQUEST")
	}

	dateFrom = time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	dateTo = time.Date(dateTo.Year(), dateTo.Month(), dateTo.Day(), 23, 59, 59, 999999999, time.UTC)

	return dateFrom, dateTo, nil
}
