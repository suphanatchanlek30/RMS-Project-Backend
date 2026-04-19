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
