package models

import "time"

type RestaurantTable struct {
	TableID     int       `json:"tableId"`
	TableNumber string    `json:"tableNumber"`
	Capacity    int       `json:"capacity"`
	TableStatus string    `json:"tableStatus"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateTableRequest struct {
	TableNumber string `json:"tableNumber"`
	Capacity    int    `json:"capacity"`
}

type UpdateTableRequest struct {
	TableNumber string `json:"tableNumber"`
	Capacity    int    `json:"capacity"`
}
