package services

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type TableService struct {
	repo *repositories.TableRepository
}

func NewTableService(repo *repositories.TableRepository) *TableService {
	return &TableService{repo: repo}
}

func (s *TableService) GetAll(ctx context.Context, status string, page, limit int) ([]models.RestaurantTable, error) {
	return s.repo.GetAll(ctx, status, page, limit)
}
