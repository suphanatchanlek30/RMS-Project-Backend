package models

type TopMenuSummary struct {
	MenuID    int    `json:"menuId"`
	MenuName  string `json:"menuName"`
	TotalSold int    `json:"totalSold"`
}

type DashboardSummaryResponse struct {
	TodaySales      float64         `json:"todaySales"`
	TodayOrders     int             `json:"todayOrders"`
	OccupiedTables  int             `json:"occupiedTables"`
	AvailableTables int             `json:"availableTables"`
	TopMenu         *TopMenuSummary `json:"topMenu"`
}
