package models

import "time"

type OrderItemRequest struct {
	MenuID   int `json:"menuId"`
	Quantity int `json:"quantity"`
}

type CreateCustomerOrderRequest struct {
	QRToken string             `json:"qrToken"`
	Items   []OrderItemRequest `json:"items"`
}

type CreateOrderRequest struct {
	SessionID           int                `json:"sessionId"`
	TableID             int                `json:"tableId"`
	CreatedByEmployeeID *int               `json:"createdByEmployeeId"`
	Items               []OrderItemRequest `json:"items"`
}

type OrderItemInput struct {
	MenuID    int
	MenuName  string
	Quantity  int
	UnitPrice float64
}

type OrderItemResponse struct {
	OrderItemID int     `json:"orderItemId"`
	MenuID      int     `json:"menuId"`
	MenuName    string  `json:"menuName"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	ItemStatus  string  `json:"itemStatus"`
}

type CreateCustomerOrderResponse struct {
	OrderID     int                 `json:"orderId"`
	SessionID   int                 `json:"sessionId"`
	TableID     int                 `json:"tableId"`
	OrderTime   time.Time           `json:"orderTime"`
	OrderStatus string              `json:"orderStatus"`
	Items       []OrderItemResponse `json:"items"`
}

type CreateOrderResponse struct {
	OrderID             int                 `json:"orderId"`
	SessionID           int                 `json:"sessionId"`
	TableID             int                 `json:"tableId"`
	CreatedByEmployeeID *int                `json:"createdByEmployeeId"`
	OrderTime           time.Time           `json:"orderTime"`
	OrderStatus         string              `json:"orderStatus"`
	Items               []OrderItemResponse `json:"items"`
}

type OrderDetailResponse struct {
	OrderID             int       `json:"orderId"`
	SessionID           int       `json:"sessionId"`
	TableID             int       `json:"tableId"`
	CreatedByEmployeeID *int      `json:"createdByEmployeeId"`
	OrderTime           time.Time `json:"orderTime"`
	OrderStatus         string    `json:"orderStatus"`
}

type CustomerOrderItemSummary struct {
	OrderItemID int     `json:"orderItemId"`
	MenuName    string  `json:"menuName"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	ItemStatus  string  `json:"itemStatus"`
}

type CustomerOrderSummary struct {
	OrderID     int                        `json:"orderId"`
	OrderTime   time.Time                  `json:"orderTime"`
	OrderStatus string                     `json:"orderStatus"`
	Items       []CustomerOrderItemSummary `json:"items"`
}

type SessionOrderSummary struct {
	OrderID     int       `json:"orderId"`
	OrderTime   time.Time `json:"orderTime"`
	OrderStatus string    `json:"orderStatus"`
}

type OrderRecord struct {
	OrderID             int
	SessionID           int
	TableID             int
	CreatedByEmployeeID *int
	OrderTime           time.Time
	OrderStatus         string
	Items               []OrderItemResponse
}

type UpdateOrderItemRequest struct {
	Quantity int `json:"quantity"`
}

type OrderItemQuantityResponse struct {
	OrderItemID int `json:"orderItemId"`
	Quantity    int `json:"quantity"`
}

type OrderItemStatusResponse struct {
	OrderItemID int    `json:"orderItemId"`
	ItemStatus  string `json:"itemStatus"`
}

type UpdateOrderItemStatusRequest struct {
	Status          string `json:"status"`
	UpdatedByChefID int    `json:"updatedByChefId"`
}

type UpdateOrderItemStatusResponse struct {
	OrderItemID int       `json:"orderItemId"`
	OldStatus   string    `json:"oldStatus"`
	NewStatus   string    `json:"newStatus"`
	UpdatedTime time.Time `json:"updatedTime"`
}

type OrderItemStatusHistory struct {
	StatusHistoryID int       `json:"statusHistoryId"`
	Status          string    `json:"status"`
	UpdatedByChefID *int      `json:"updatedByChefId"`
	UpdatedTime     time.Time `json:"updatedTime"`
}
