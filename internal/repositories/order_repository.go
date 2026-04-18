package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
)

type OrderRepository struct {
	DB *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, sessionID int, tableID int, createdByEmployeeID *int, items []models.OrderItemInput) (*models.OrderRecord, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO customer_orders (session_id, table_id, created_by_employee_id)
		VALUES ($1, $2, $3)
		RETURNING order_id, session_id, table_id, order_time, order_status
	`

	var record models.OrderRecord
	err = tx.QueryRow(ctx, query, sessionID, tableID, createdByEmployeeID).Scan(
		&record.OrderID,
		&record.SessionID,
		&record.TableID,
		&record.OrderTime,
		&record.OrderStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	record.CreatedByEmployeeID = createdByEmployeeID

	for _, item := range items {
		itemQuery := `
			INSERT INTO order_items (order_id, menu_id, quantity, unit_price, item_status)
			VALUES ($1, $2, $3, $4, 'WAITING')
			RETURNING order_item_id, item_status
		`

		var orderItem models.OrderItemResponse
		err = tx.QueryRow(ctx, itemQuery, record.OrderID, item.MenuID, item.Quantity, item.UnitPrice).Scan(
			&orderItem.OrderItemID,
			&orderItem.ItemStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}

		orderItem.MenuID = item.MenuID
		orderItem.MenuName = item.MenuName
		orderItem.Quantity = item.Quantity
		orderItem.UnitPrice = item.UnitPrice
		record.Items = append(record.Items, orderItem)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return &record, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, orderID int) (*models.OrderRecord, error) {
	query := `
		SELECT
			co.order_id,
			co.session_id,
			co.table_id,
			co.created_by_employee_id,
			co.order_time,
			co.order_status,
			oi.order_item_id,
			oi.menu_id,
			m.menu_name,
			oi.quantity,
			oi.unit_price,
			oi.item_status
		FROM customer_orders co
		JOIN order_items oi ON co.order_id = oi.order_id
		JOIN menus m ON oi.menu_id = m.menu_id
		WHERE co.order_id = $1
		ORDER BY oi.order_item_id ASC
	`

	rows, err := r.DB.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	var record *models.OrderRecord
	for rows.Next() {
		var createdBy sql.NullInt32
		var item models.OrderItemResponse

		if record == nil {
			record = &models.OrderRecord{}
		}

		err := rows.Scan(
			&record.OrderID,
			&record.SessionID,
			&record.TableID,
			&createdBy,
			&record.OrderTime,
			&record.OrderStatus,
			&item.OrderItemID,
			&item.MenuID,
			&item.MenuName,
			&item.Quantity,
			&item.UnitPrice,
			&item.ItemStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}

		if createdBy.Valid {
			value := int(createdBy.Int32)
			record.CreatedByEmployeeID = &value
		}
		record.Items = append(record.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if record == nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	return record, nil
}

func (r *OrderRepository) GetBySessionID(ctx context.Context, sessionID int) ([]models.OrderRecord, error) {
	query := `
		SELECT
			co.order_id,
			co.session_id,
			co.table_id,
			co.created_by_employee_id,
			co.order_time,
			co.order_status,
			oi.order_item_id,
			oi.menu_id,
			m.menu_name,
			oi.quantity,
			oi.unit_price,
			oi.item_status
		FROM customer_orders co
		JOIN order_items oi ON co.order_id = oi.order_id
		JOIN menus m ON oi.menu_id = m.menu_id
		WHERE co.session_id = $1
		ORDER BY co.order_time ASC, co.order_id ASC, oi.order_item_id ASC
	`

	rows, err := r.DB.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	defer rows.Close()

	var records []models.OrderRecord
	indexByOrderID := make(map[int]int)

	for rows.Next() {
		var createdBy sql.NullInt32
		var item models.OrderItemResponse
		var order models.OrderRecord

		err := rows.Scan(
			&order.OrderID,
			&order.SessionID,
			&order.TableID,
			&createdBy,
			&order.OrderTime,
			&order.OrderStatus,
			&item.OrderItemID,
			&item.MenuID,
			&item.MenuName,
			&item.Quantity,
			&item.UnitPrice,
			&item.ItemStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("INTERNAL")
		}

		if idx, ok := indexByOrderID[order.OrderID]; ok {
			records[idx].Items = append(records[idx].Items, item)
			continue
		}

		if createdBy.Valid {
			value := int(createdBy.Int32)
			order.CreatedByEmployeeID = &value
		}
		order.Items = []models.OrderItemResponse{item}
		indexByOrderID[order.OrderID] = len(records)
		records = append(records, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	if records == nil {
		records = []models.OrderRecord{}
	}

	return records, nil
}

func (r *OrderRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItemResponse, error) {
	query := `
		SELECT 
			oi.order_item_id,
			oi.menu_id,
			m.menu_name,
			oi.quantity,
			oi.unit_price,
			oi.item_status
		FROM order_items oi
		JOIN menus m ON oi.menu_id = m.menu_id
		WHERE oi.order_id = $1;
	`

	rows, err := r.DB.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItemResponse

	for rows.Next() {
		var item models.OrderItemResponse
		err := rows.Scan(
			&item.OrderItemID,
			&item.MenuID,
			&item.MenuName,
			&item.Quantity,
			&item.UnitPrice,
			&item.ItemStatus,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, pgx.ErrNoRows
	}

	return items, nil
}

func (r *OrderRepository) GetOrderItemByID(ctx context.Context, id int) (*models.OrderItemResponse, error) {
	query := `
		SELECT order_item_id, quantity, item_status
		FROM order_items
		WHERE order_item_id = $1
	`

	var item models.OrderItemResponse

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&item.OrderItemID,
		&item.Quantity,
		&item.ItemStatus,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r *OrderRepository) UpdateOrderItemQuantity(ctx context.Context, id int, quantity int) (*models.OrderItemQuantityResponse, error) {
	query := `
		UPDATE order_items
		SET quantity = $1
		WHERE order_item_id = $2
		RETURNING order_item_id, quantity
	`

	var res models.OrderItemQuantityResponse

	err := r.DB.QueryRow(ctx, query, quantity, id).Scan(
		&res.OrderItemID,
		&res.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *OrderRepository) UpdateOrderItemStatus(ctx context.Context, id int, status string) (*models.OrderItemStatusResponse, error) {
	query := `
		UPDATE order_items
		SET item_status = $1
		WHERE order_item_id = $2
		RETURNING order_item_id, item_status
	`

	var res models.OrderItemStatusResponse

	err := r.DB.QueryRow(ctx, query, status, id).Scan(
		&res.OrderItemID,
		&res.ItemStatus,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
