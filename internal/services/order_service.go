package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"time"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type OrderService struct {
	repo             *repositories.OrderRepository
	qrSessionRepo    *repositories.QRSessionRepository
	tableSessionRepo *repositories.TableSessionRepository
	tableRepo        *repositories.TableRepository
	menuRepo         *repositories.MenuRepository
}

func NewOrderService(repo *repositories.OrderRepository, qrSessionRepo *repositories.QRSessionRepository, tableSessionRepo *repositories.TableSessionRepository, tableRepo *repositories.TableRepository, menuRepo *repositories.MenuRepository) *OrderService {
	return &OrderService{repo: repo, qrSessionRepo: qrSessionRepo, tableSessionRepo: tableSessionRepo, tableRepo: tableRepo, menuRepo: menuRepo}
}

func (s *OrderService) CreateCustomerOrder(ctx context.Context, req models.CreateCustomerOrderRequest) (*models.CreateCustomerOrderResponse, error) {
	sessionID, tableID, err := s.validateCustomerContext(ctx, req.QRToken)
	if err != nil {
		return nil, err
	}

	validated, err := s.validateMenus(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	order, err := s.repo.CreateOrder(ctx, sessionID, tableID, nil, validated)
	if err != nil {
		return nil, err
	}

	return toCustomerOrderCreateResponse(order), nil
}

func (s *OrderService) GetCustomerOrders(ctx context.Context, qrToken string) ([]models.CustomerOrderSummary, error) {
	sessionID, _, err := s.validateCustomerContext(ctx, qrToken)
	if err != nil {
		return nil, err
	}

	orders, err := s.repo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return toCustomerOrderSummaries(orders), nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req models.CreateOrderRequest) (*models.CreateOrderResponse, error) {
	if req.SessionID <= 0 || req.TableID <= 0 || req.CreatedByEmployeeID == nil || *req.CreatedByEmployeeID <= 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	session, err := s.tableSessionRepo.GetSessionByID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	if session.SessionStatus != "OPEN" {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	if session.TableID != req.TableID {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	table, err := s.tableRepo.GetByID(ctx, req.TableID)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	if table.TableStatus != "OCCUPIED" {
		return nil, fmt.Errorf("UNPROCESSABLE")
	}

	validated, err := s.validateMenus(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	order, err := s.repo.CreateOrder(ctx, req.SessionID, req.TableID, req.CreatedByEmployeeID, validated)
	if err != nil {
		return nil, err
	}

	return toCashierOrderCreateResponse(order), nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID int) (*models.OrderDetailResponse, error) {
	if orderID <= 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return &models.OrderDetailResponse{
		OrderID:             order.OrderID,
		SessionID:           order.SessionID,
		TableID:             order.TableID,
		CreatedByEmployeeID: order.CreatedByEmployeeID,
		OrderTime:           order.OrderTime,
		OrderStatus:         order.OrderStatus,
	}, nil
}

func (s *OrderService) GetOrdersBySessionID(ctx context.Context, sessionID int) ([]models.SessionOrderSummary, error) {
	if sessionID <= 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	if _, err := s.tableSessionRepo.GetSessionByID(ctx, sessionID); err != nil {
		return nil, fmt.Errorf("NOT_FOUND")
	}

	orders, err := s.repo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return toSessionOrderSummaries(orders), nil
}

func (s *OrderService) validateCustomerContext(ctx context.Context, qrToken string) (int, int, error) {
	if qrToken == "" {
		return 0, 0, fmt.Errorf("BAD_REQUEST")
	}

	qr, err := s.qrSessionRepo.GetByToken(ctx, qrToken)
	if err != nil {
		return 0, 0, fmt.Errorf("NOT_FOUND")
	}

	if time.Now().After(qr.ExpiredAt) {
		return 0, 0, fmt.Errorf("GONE")
	}

	if qr.SessionStatus == "CLOSED" {
		return 0, 0, fmt.Errorf("UNPROCESSABLE")
	}

	table, err := s.tableRepo.GetByID(ctx, qr.TableID)
	if err != nil {
		return 0, 0, fmt.Errorf("NOT_FOUND")
	}

	if table.TableStatus != "OCCUPIED" {
		return 0, 0, fmt.Errorf("UNPROCESSABLE")
	}

	return qr.SessionID, qr.TableID, nil
}

func (s *OrderService) validateMenus(ctx context.Context, items []models.OrderItemRequest) ([]models.OrderItemInput, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	validated := make([]models.OrderItemInput, 0, len(items))
	for _, item := range items {
		if item.MenuID <= 0 || item.Quantity <= 0 {
			return nil, fmt.Errorf("BAD_REQUEST")
		}

		menu, err := s.menuRepo.GetByID(ctx, item.MenuID)
		if err != nil {
			if err.Error() == "NOT_FOUND" {
				return nil, fmt.Errorf("NOT_FOUND")
			}
			return nil, fmt.Errorf("INTERNAL")
		}

		if !menu.MenuStatus {
			return nil, fmt.Errorf("UNPROCESSABLE")
		}

		validated = append(validated, models.OrderItemInput{
			MenuID:    menu.MenuID,
			MenuName:  menu.MenuName,
			Quantity:  item.Quantity,
			UnitPrice: menu.Price,
		})
	}

	return validated, nil
}

func toCustomerOrderCreateResponse(order *models.OrderRecord) *models.CreateCustomerOrderResponse {
	items := make([]models.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, item)
	}

	return &models.CreateCustomerOrderResponse{
		OrderID:     order.OrderID,
		SessionID:   order.SessionID,
		TableID:     order.TableID,
		OrderTime:   order.OrderTime,
		OrderStatus: order.OrderStatus,
		Items:       items,
	}
}

func toCashierOrderCreateResponse(order *models.OrderRecord) *models.CreateOrderResponse {
	items := make([]models.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, item)
	}

	return &models.CreateOrderResponse{
		OrderID:             order.OrderID,
		SessionID:           order.SessionID,
		TableID:             order.TableID,
		CreatedByEmployeeID: order.CreatedByEmployeeID,
		OrderTime:           order.OrderTime,
		OrderStatus:         order.OrderStatus,
		Items:               items,
	}
}

func toCustomerOrderSummaries(records []models.OrderRecord) []models.CustomerOrderSummary {
	result := make([]models.CustomerOrderSummary, 0, len(records))
	for _, record := range records {
		items := make([]models.CustomerOrderItemSummary, 0, len(record.Items))
		for _, item := range record.Items {
			items = append(items, models.CustomerOrderItemSummary{
				OrderItemID: item.OrderItemID,
				MenuName:    item.MenuName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				ItemStatus:  item.ItemStatus,
			})
		}

		result = append(result, models.CustomerOrderSummary{
			OrderID:     record.OrderID,
			OrderTime:   record.OrderTime,
			OrderStatus: record.OrderStatus,
			Items:       items,
		})
	}

	if result == nil {
		result = []models.CustomerOrderSummary{}
	}

	return result
}

func toSessionOrderSummaries(records []models.OrderRecord) []models.SessionOrderSummary {
	result := make([]models.SessionOrderSummary, 0, len(records))
	for _, record := range records {
		result = append(result, models.SessionOrderSummary{
			OrderID:     record.OrderID,
			OrderTime:   record.OrderTime,
			OrderStatus: record.OrderStatus,
		})
	}

	if result == nil {
		result = []models.SessionOrderSummary{}
	}

	return result
}

func (s *OrderService) GetOrderItems(ctx context.Context, orderID int) ([]models.OrderItemResponse, error) {
	items, err := s.repo.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return items, nil
}
