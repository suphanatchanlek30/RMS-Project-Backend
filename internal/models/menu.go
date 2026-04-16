package models

import "time"

type Menu struct {
	MenuID       int       `json:"menuId"`
	MenuName     string    `json:"menuName"`
	CategoryID   int       `json:"categoryId"`
	CategoryName string    `json:"categoryName,omitempty"`
	Price        float64   `json:"price"`
	Description  string    `json:"description"`
	MenuStatus   bool      `json:"menuStatus"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CreateMenuRequest struct {
	MenuName    string  `json:"menuName"`
	CategoryID  int     `json:"categoryId"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	MenuStatus  bool    `json:"menuStatus"`
}

type CreateMenuResponse struct {
	MenuID      int       `json:"menuId"`
	MenuName    string    `json:"menuName"`
	CategoryID  int       `json:"categoryId"`
	Price       float64   `json:"price"`
	Description *string   `json:"description"`
	MenuStatus  bool      `json:"menuStatus"`
	CreatedAt   time.Time `json:"createdAt"`
}
