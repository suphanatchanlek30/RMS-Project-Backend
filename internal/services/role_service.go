package services

import (
	"context"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type RoleService struct {
	repo *repositories.RoleRepository
}

func NewRoleService(repo *repositories.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAll(ctx context.Context) ([]models.Role, error) {
	return s.repo.GetAll(ctx)
}