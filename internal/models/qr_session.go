package models

import "time"

type CreateQRSessionRequest struct {
	SessionID int `json:"sessionId"`
}

type CreateQRSessionResponse struct {
	QRSessionID int       `json:"qrSessionId"`
	SessionID   int       `json:"sessionId"`
	QRCodeURL   string    `json:"qrCodeUrl"`
	QRToken     string    `json:"qrToken"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiredAt   time.Time `json:"expiredAt"`
}
