package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/config"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type QRSessionRepository struct {
	DB *pgxpool.Pool
}

func NewQRSessionRepository(db *pgxpool.Pool) *QRSessionRepository {
	return &QRSessionRepository{DB: db}
}

func (r *QRSessionRepository) HasActiveQR(ctx context.Context, sessionID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM qr_sessions
			WHERE session_id = $1 AND expired_at > CURRENT_TIMESTAMP
		)
	`

	var exists bool
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *QRSessionRepository) CreateQRSession(ctx context.Context, sessionID int) (*models.CreateQRSessionResponse, error) {
	qrToken := strings.ReplaceAll(uuid.New().String(), "-", "")
	frontendURL := config.GetEnv("FRONTEND_URL", "http://localhost:3000")
	qrCodeURL := frontendURL + "/q/" + qrToken
	expiredAt := time.Now().Add(4 * time.Hour)

	query := `
		INSERT INTO qr_sessions (session_id, qr_code_url, expired_at)
		VALUES ($1, $2, $3)
		RETURNING qr_session_id, session_id, qr_code_url, created_at, expired_at
	`

	var resp models.CreateQRSessionResponse
	err := r.DB.QueryRow(ctx, query, sessionID, qrCodeURL, expiredAt).Scan(
		&resp.QRSessionID,
		&resp.SessionID,
		&resp.QRCodeURL,
		&resp.CreatedAt,
		&resp.ExpiredAt,
	)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	resp.QRToken = qrToken

	return &resp, nil
}
