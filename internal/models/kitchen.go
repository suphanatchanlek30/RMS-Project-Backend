package models

import "time"

type KitchenOrderResponse struct {
	OrderID     int           `json:"orderId"`
	TableID     int           `json:"tableId"`
	TableNumber string        `json:"tableNumber"`
	OrderTime   time.Time     `json:"orderTime"`
	Items       []KitchenItem `json:"items"`
}

type KitchenItem struct {
	OrderItemID int    `json:"orderItemId"`
	MenuName    string `json:"menuName"`
	Quantity    int    `json:"quantity"`
	ItemStatus  string `json:"itemStatus"`
}
