package models

import "time"

type CashierCurrentSession struct {
	SessionID int       `json:"sessionId"`
	StartTime time.Time `json:"startTime"`
}

type CashierTableOverview struct {
	TableID        int                    `json:"tableId"`
	TableNumber    string                 `json:"tableNumber"`
	TableStatus    string                 `json:"tableStatus"`
	CurrentSession *CashierCurrentSession `json:"currentSession"`
}

type CheckoutBillItem struct {
	OrderItemID int     `json:"orderItemId"`
	MenuName    string  `json:"menuName"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	LineTotal    float64 `json:"lineTotal"`
}

type CheckoutBill struct {
	Items       []CheckoutBillItem `json:"items"`
	TotalAmount float64            `json:"totalAmount"`
}

type CheckoutPaymentMethod struct {
	PaymentMethodID int    `json:"paymentMethodId"`
	MethodName      string `json:"methodName"`
}

type CheckoutResponse struct {
	SessionID      int                     `json:"sessionId"`
	TableID        int                     `json:"tableId"`
	TableNumber    string                  `json:"tableNumber"`
	Bill           CheckoutBill             `json:"bill"`
	PaymentMethods []CheckoutPaymentMethod `json:"paymentMethods"`
}