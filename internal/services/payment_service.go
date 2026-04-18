package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type PaymentService struct {
	paymentRepo      *repositories.PaymentRepository
	tableSessionRepo *repositories.TableSessionRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, tableSessionRepo *repositories.TableSessionRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, tableSessionRepo: tableSessionRepo}
}

func (s *PaymentService) Create(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	if req.SessionID <= 0 || req.PaymentMethodID <= 0 || req.ReceivedAmount < 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	if _, err := s.tableSessionRepo.GetSessionByID(ctx, req.SessionID); err != nil {
		return nil, fmt.Errorf("NOT_FOUND_SESSION")
	}

	methodName, err := s.paymentRepo.GetPaymentMethodNameByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("NOT_FOUND_PAYMENT_METHOD")
	}

	hasPaid, err := s.paymentRepo.HasPaidPaymentBySession(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasPaid {
		return nil, fmt.Errorf("CONFLICT")
	}

	hasUnbillable, err := s.tableSessionRepo.HasUnbillableItems(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}
	if hasUnbillable {
		return nil, fmt.Errorf("NOT_READY")
	}

	bill, err := s.tableSessionRepo.GetSessionBill(ctx, req.SessionID)
	if err != nil {
		if err.Error() == "NO_ITEMS" {
			return nil, fmt.Errorf("NOT_READY")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	if req.ReceivedAmount < bill.TotalAmount {
		return nil, fmt.Errorf("INSUFFICIENT")
	}

	orderID, err := s.paymentRepo.GetLatestOrderIDBySession(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("NOT_READY")
	}

	resp, err := s.paymentRepo.CreatePayment(
		ctx,
		orderID,
		req.SessionID,
		req.PaymentMethodID,
		methodName,
		bill.TotalAmount,
		req.ReceivedAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}

func (s *PaymentService) GetByID(ctx context.Context, paymentID int) (*models.PaymentDetailResponse, error) {
	if paymentID <= 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	resp, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}

func (s *PaymentService) GetAll(ctx context.Context, filter models.PaymentListFilter) ([]models.PaymentListItem, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	items, total, err := s.paymentRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("INTERNAL")
	}

	return items, total, nil
}
