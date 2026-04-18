package services

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type KitchenService struct {
	repo *repositories.KitchenRepository
}

func NewKitchenService(r *repositories.KitchenRepository) *KitchenService {
	return &KitchenService{repo: r}
}

func (s *KitchenService) GetKitchenOrders(
	ctx context.Context,
	status string,
	tableID int,
	page int,
	limit int,
) ([]models.KitchenOrderResponse, error) {

	offset := (page - 1) * limit

	return s.repo.GetKitchenOrders(ctx, status, tableID, limit, offset)
}
