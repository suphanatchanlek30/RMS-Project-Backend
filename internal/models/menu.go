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

type MenuDetail struct {
	MenuID      int     `json:"menuId"`
	MenuName    string  `json:"menuName"`
	CategoryID  int     `json:"categoryId"`
	Price       float64 `json:"price"`
	Description *string `json:"description"`
	MenuStatus  bool    `json:"menuStatus"`
}

type UpdateMenuRequest struct {
	MenuName    string  `json:"menuName"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type UpdateMenuResponse struct {
	MenuID      int     `json:"menuId"`
	MenuName    string  `json:"menuName"`
	Price       float64 `json:"price"`
	Description *string `json:"description"`
}

type UpdateMenuStatusRequest struct {
	MenuStatus *bool `json:"menuStatus"`
}

type UpdateMenuStatusResponse struct {
	MenuID     int  `json:"menuId"`
	MenuStatus bool `json:"menuStatus"`
}

type CustomerMenuTable struct {
	TableID     int    `json:"tableId"`
	TableNumber string `json:"tableNumber"`
}

type CustomerMenuCategory struct {
	CategoryID   int    `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}

type CustomerMenuItem struct {
	MenuID      int     `json:"menuId"`
	MenuName    string  `json:"menuName"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	MenuStatus  bool    `json:"menuStatus"`
}

type CustomerMenuResponse struct {
	Table      CustomerMenuTable      `json:"table"`
	Categories []CustomerMenuCategory `json:"categories"`
	Menus      []CustomerMenuItem     `json:"menus"`
}
