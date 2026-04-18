package services

import (
	"context"
	"fmt"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
)

type ReceiptService struct {
	repo        *repositories.ReceiptRepository
	paymentRepo *repositories.PaymentRepository
}

func NewReceiptService(repo *repositories.ReceiptRepository, paymentRepo *repositories.PaymentRepository) *ReceiptService {
	return &ReceiptService{repo: repo, paymentRepo: paymentRepo}
}

func (s *ReceiptService) GetByPaymentID(ctx context.Context, paymentID int) (*models.ReceiptDetailResponse, error) {
	if paymentID <= 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	resp, err := s.repo.GetByPaymentID(ctx, paymentID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			payment, perr := s.paymentRepo.GetByID(ctx, paymentID)
			if perr != nil {
				if perr.Error() == "NOT_FOUND" {
					return nil, fmt.Errorf("NOT_FOUND")
				}
				return nil, fmt.Errorf("INTERNAL")
			}

			if cerr := s.repo.CreateForPayment(ctx, paymentID, payment.TotalAmount); cerr != nil {
				return nil, fmt.Errorf("INTERNAL")
			}

			resp, rerr := s.repo.GetByPaymentID(ctx, paymentID)
			if rerr != nil {
				if rerr.Error() == "NOT_FOUND" {
					return nil, fmt.Errorf("NOT_FOUND")
				}
				return nil, fmt.Errorf("INTERNAL")
			}

			return resp, nil
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}

func (s *ReceiptService) GetByReceiptID(ctx context.Context, receiptID int) (*models.ReceiptDetailResponse, error) {
	if receiptID <= 0 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}

	resp, err := s.repo.GetByReceiptID(ctx, receiptID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		return nil, fmt.Errorf("INTERNAL")
	}

	return resp, nil
}
