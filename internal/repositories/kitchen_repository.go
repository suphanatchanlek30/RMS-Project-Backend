package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type KitchenRepository struct {
	DB *pgxpool.Pool
}

func NewKitchenRepository(db *pgxpool.Pool) *KitchenRepository {
	return &KitchenRepository{DB: db}
}

func (r *KitchenRepository) GetKitchenOrders(
	ctx context.Context,
	status string,
	tableID int,
	limit int,
	offset int,
) ([]models.KitchenOrderResponse, error) {

	query := `
	SELECT 
		o.order_id,
		t.table_id,
		t.table_number,
		o.order_time,
		oi.order_item_id,
		m.menu_name,
		oi.quantity,
		oi.item_status
	FROM customer_orders o
	JOIN table_sessions ts ON o.session_id = ts.session_id
	JOIN restaurant_tables t ON ts.table_id = t.table_id
	JOIN order_items oi ON o.order_id = oi.order_id
	JOIN menus m ON oi.menu_id = m.menu_id
	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	// 🔹 filter status
	if status != "" {
		query += fmt.Sprintf(" AND oi.item_status = $%d", argPos)
		args = append(args, status)
		argPos++
	}

	// 🔹 filter table
	if tableID != 0 {
		query += fmt.Sprintf(" AND t.table_id = $%d", argPos)
		args = append(args, tableID)
		argPos++
	}

	// 🔹 default kitchen filter (ถ้าไม่ส่ง status)
	if status == "" {
		query += " AND oi.item_status IN ('WAITING','PREPARING')"
	}

	// 🔹 pagination
	query += fmt.Sprintf(" ORDER BY o.order_time ASC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[int]*models.KitchenOrderResponse)

	for rows.Next() {
		var orderID, tableID int
		var tableNumber, menuName, itemStatus string
		var orderTime time.Time
		var orderItemID, quantity int

		err := rows.Scan(
			&orderID,
			&tableID,
			&tableNumber,
			&orderTime,
			&orderItemID,
			&menuName,
			&quantity,
			&itemStatus,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := orderMap[orderID]; !ok {
			orderMap[orderID] = &models.KitchenOrderResponse{
				OrderID:     orderID,
				TableID:     tableID,
				TableNumber: tableNumber,
				OrderTime:   orderTime,
				Items:       []models.KitchenItem{},
			}
		}

		orderMap[orderID].Items = append(orderMap[orderID].Items, models.KitchenItem{
			OrderItemID: orderItemID,
			MenuName:    menuName,
			Quantity:    quantity,
			ItemStatus:  itemStatus,
		})
	}

	var result []models.KitchenOrderResponse
	for _, v := range orderMap {
		result = append(result, *v)
	}

	return result, nil
}
