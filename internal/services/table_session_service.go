package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type TableSessionService struct {
	repo *repositories.TableSessionRepository
}

func NewTableSessionService(repo *repositories.TableSessionRepository) *TableSessionService {
	return &TableSessionService{repo: repo}
}

func (s *TableSessionService) OpenTable(ctx context.Context, req models.OpenTableRequest) (*models.OpenTableResponse, error) {
	table, err := s.repo.GetTableByID(ctx, req.TableID)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	if table.TableStatus == "OCCUPIED" {
		return nil, fmt.Errorf("CONFLICT")
	}

	if table.TableStatus != "AVAILABLE" {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	hasOpen, err := s.repo.HasOpenSession(ctx, req.TableID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasOpen {
		return nil, fmt.Errorf("CONFLICT")
	}

	resp, err := s.repo.OpenSession(ctx, req.TableID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}

func (s *TableSessionService) GetByID(ctx context.Context, sessionID int) (*models.TableSessionDetail, error) {
	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	return session, nil
}
