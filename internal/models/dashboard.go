package models

type TopMenuInfo struct {
	MenuID    int    `json:"menuId"`
	MenuName  string `json:"menuName"`
	TotalSold int    `json:"totalSold"`
}

type DashboardSummaryData struct {
	TodaySales     float64      `json:"todaySales"`
	TodayOrders    int          `json:"todayOrders"`
	OccupiedTables int          `json:"occupiedTables"`
	AvailableTables int         `json:"availableTables"`
	TopMenu        TopMenuInfo  `json:"topMenu"`
}
