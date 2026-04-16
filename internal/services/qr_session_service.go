package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type QRSessionService struct {
	repo            *repositories.QRSessionRepository
	tableSessionRepo *repositories.TableSessionRepository
}

func NewQRSessionService(repo *repositories.QRSessionRepository, tableSessionRepo *repositories.TableSessionRepository) *QRSessionService {
	return &QRSessionService{repo: repo, tableSessionRepo: tableSessionRepo}
}

func (s *QRSessionService) CreateQRSession(ctx context.Context, req models.CreateQRSessionRequest) (*models.CreateQRSessionResponse, error) {
	session, err := s.tableSessionRepo.GetSessionByID(ctx, req.SessionID)
	if err != nil || session == nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	if session.SessionStatus != "OPEN" {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	hasActive, err := s.repo.HasActiveQR(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasActive {
		return nil, fmt.Errorf("CONFLICT")
	}

	resp, err := s.repo.CreateQRSession(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}
