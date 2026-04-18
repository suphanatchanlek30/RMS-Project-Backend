package models

import "time"

type CreatePaymentRequest struct {
	SessionID       int     `json:"sessionId"`
	PaymentMethodID int     `json:"paymentMethodId"`
	ReceivedAmount  float64 `json:"receivedAmount"`
}

type CreatePaymentResponse struct {
	PaymentID         int       `json:"paymentId"`
	SessionID         int       `json:"sessionId"`
	PaymentMethodID   int       `json:"paymentMethodId"`
	PaymentMethodName string    `json:"paymentMethodName"`
	TotalAmount       float64   `json:"totalAmount"`
	ReceivedAmount    float64   `json:"receivedAmount"`
	ChangeAmount      float64   `json:"changeAmount"`
	PaymentTime       time.Time `json:"paymentTime"`
	PaymentStatus     string    `json:"paymentStatus"`
}

type PaymentDetailResponse struct {
	PaymentID         int       `json:"paymentId"`
	SessionID         int       `json:"sessionId"`
	PaymentMethodID   int       `json:"paymentMethodId"`
	PaymentMethodName string    `json:"paymentMethodName"`
	TotalAmount       float64   `json:"totalAmount"`
	PaymentTime       time.Time `json:"paymentTime"`
	PaymentStatus     string    `json:"paymentStatus"`
}

type PaymentListItem struct {
	PaymentID         int       `json:"paymentId"`
	SessionID         int       `json:"sessionId"`
	PaymentMethodName string    `json:"paymentMethodName"`
	TotalAmount       float64   `json:"totalAmount"`
	PaymentStatus     string    `json:"paymentStatus"`
	PaymentTime       time.Time `json:"paymentTime"`
}

type PaymentListFilter struct {
	DateFrom        string
	DateTo          string
	PaymentMethodID *int
	Status          string
	Page            int
	Limit           int
}
