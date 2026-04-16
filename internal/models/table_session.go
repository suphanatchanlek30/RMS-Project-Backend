package models

import "time"

type TableSession struct {
	SessionID      int       `json:"sessionId"`
	TableID        int       `json:"tableId"`
	StartTime      time.Time `json:"startTime"`
	EndTime        *time.Time `json:"endTime,omitempty"`
	SessionDetails *string   `json:"sessionDetails,omitempty"`
	SessionStatus  string    `json:"sessionStatus"`
}

type OpenTableRequest struct {
	TableID    int `json:"tableId"`
	EmployeeID int `json:"employeeId"`
}

type OpenTableResponse struct {
	SessionID     int       `json:"sessionId"`
	TableID       int       `json:"tableId"`
	TableNumber   string    `json:"tableNumber"`
	StartTime     time.Time `json:"startTime"`
	SessionStatus string    `json:"sessionStatus"`
}

type TableSessionDetail struct {
	SessionID     int        `json:"sessionId"`
	TableID       int        `json:"tableId"`
	TableNumber   string     `json:"tableNumber"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	SessionStatus string     `json:"sessionStatus"`
}

type CurrentSessionResponse struct {
	SessionID     int       `json:"sessionId"`
	TableID       int       `json:"tableId"`
	SessionStatus string    `json:"sessionStatus"`
	StartTime     time.Time `json:"startTime"`
}
