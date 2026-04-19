package models

import "time"

type SalesReportQuery struct {
	DateFrom time.Time
	DateTo   time.Time
	GroupBy  string
}

type TopMenusReportQuery struct {
	DateFrom time.Time
	DateTo   time.Time
	Limit    int
}

type SalesReportItem struct {
	Date        string  `json:"date"`
	TotalSales  float64 `json:"totalSales"`
	TotalOrders int     `json:"totalOrders"`
}

type TopMenuReportItem struct {
	MenuID        int     `json:"menuId"`
	MenuName      string  `json:"menuName"`
	TotalQuantity int     `json:"totalQuantity"`
	TotalAmount   float64 `json:"totalAmount"`
}
