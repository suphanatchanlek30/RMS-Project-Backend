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

type VerifyQRResponse struct {
	QRSessionID   int       `json:"qrSessionId"`
	SessionID     int       `json:"sessionId"`
	TableID       int       `json:"tableId"`
	TableNumber   string    `json:"tableNumber"`
	SessionStatus string    `json:"sessionStatus"`
	ExpiredAt     time.Time `json:"expiredAt"`
}
