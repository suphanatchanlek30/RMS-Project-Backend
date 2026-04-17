package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, req models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	resp, err := s.repo.Create(ctx, req)
	if err != nil {
		if err.Error() == "CONFLICT" {
			return nil, fmt.Errorf("CONFLICT")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}

func (s *CategoryService) GetAll(ctx context.Context) ([]models.CategoryListItem, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return categories, nil
}

func (s *CategoryService) Update(ctx context.Context, categoryID int, req models.UpdateCategoryRequest) (*models.CategoryListItem, error) {
	resp, err := s.repo.Update(ctx, categoryID, req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return nil, fmt.Errorf("NOT_FOUND")
		case "CONFLICT":
			return nil, fmt.Errorf("CONFLICT")
		default:
			return nil, fmt.Errorf("INTERNAL")
		}
	}

	return resp, nil
}
