package models

import "time"

type CashierTableCurrentSession struct {
	SessionID int       `json:"sessionId"`
	StartTime time.Time `json:"startTime"`
}

type CashierTableOverviewItem struct {
	TableID        int                         `json:"tableId"`
	TableNumber    string                      `json:"tableNumber"`
	TableStatus    string                      `json:"tableStatus"`
	CurrentSession *CashierTableCurrentSession `json:"currentSession"`
}

type CheckoutBill struct {
	Items       []BillItemResponse `json:"items"`
	TotalAmount float64            `json:"totalAmount"`
}

type SessionCheckoutResponse struct {
	SessionID      int                 `json:"sessionId"`
	TableID        int                 `json:"tableId"`
	TableNumber    string              `json:"tableNumber"`
	Bill           CheckoutBill        `json:"bill"`
	PaymentMethods []PaymentMethodItem `json:"paymentMethods"`
}

type CashierCheckoutRequest struct {
	SessionID       int     `json:"sessionId"`
	PaymentMethodID int     `json:"paymentMethodId"`
	ReceivedAmount  float64 `json:"receivedAmount"`
}

type CashierCheckoutResult struct {
	PaymentID     int     `json:"paymentId"`
	ReceiptID     int     `json:"receiptId"`
	ReceiptNumber string  `json:"receiptNumber"`
	SessionID     int     `json:"sessionId"`
	SessionStatus string  `json:"sessionStatus"`
	TableID       int     `json:"tableId"`
	TableStatus   string  `json:"tableStatus"`
	ChangeAmount  float64 `json:"changeAmount"`
}
