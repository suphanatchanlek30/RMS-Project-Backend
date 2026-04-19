package models

import "time"

type SalesReportItem struct {
	Date        string  `json:"date"`
	TotalSales  float64 `json:"totalSales"`
	TotalOrders int     `json:"totalOrders"`
}

type SalesReportRequest struct {
	DateFrom string `query:"dateFrom"`
	DateTo   string `query:"dateTo"`
	GroupBy  string `query:"groupBy"`
}

type TopMenuReportItem struct {
	MenuID        int     `json:"menuId"`
	MenuName      string  `json:"menuName"`
	TotalQuantity int     `json:"totalQuantity"`
	TotalAmount   float64 `json:"totalAmount"`
}
