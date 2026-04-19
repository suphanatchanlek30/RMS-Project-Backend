package models

import "time"

type CashierTableOverviewItem struct {
	TableID        int                    `json:"tableId"`
	TableNumber    string                 `json:"tableNumber"`
	TableStatus    string                 `json:"tableStatus"`
	CurrentSession *CashierCurrentSession `json:"currentSession"`
}

type CashierCurrentSession struct {
	SessionID int       `json:"sessionId"`
	StartTime time.Time `json:"startTime"`
}

type CashierTablesOverviewResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    []CashierTableOverviewItem `json:"data"`
}