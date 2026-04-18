package models

import "time"

type ReceiptPaymentInfo struct {
	PaymentID         int       `json:"paymentId"`
	PaymentMethodName string    `json:"paymentMethodName"`
	PaymentTime       time.Time `json:"paymentTime"`
}

type ReceiptTableInfo struct {
	TableID     int    `json:"tableId"`
	TableNumber string `json:"tableNumber"`
}

type ReceiptItem struct {
	MenuName  string  `json:"menuName"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	LineTotal float64 `json:"lineTotal"`
}

type ReceiptDetailResponse struct {
	ReceiptID     int                `json:"receiptId"`
	ReceiptNumber string             `json:"receiptNumber"`
	IssueDate     time.Time          `json:"issueDate"`
	TotalAmount   float64            `json:"totalAmount"`
	Payment       ReceiptPaymentInfo `json:"payment"`
	Table         ReceiptTableInfo   `json:"table"`
	Items         []ReceiptItem      `json:"items"`
}
