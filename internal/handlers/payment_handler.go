package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Create(c *fiber.Ctx) error {
	var req models.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.SessionID <= 0 || req.PaymentMethodID <= 0 || req.ReceivedAmount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.Create(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "ข้อมูลไม่ถูกต้อง",
				Data:    nil,
			})
		case "NOT_FOUND_SESSION":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ session",
				Data:    nil,
			})
		case "NOT_FOUND_PAYMENT_METHOD":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ payment method",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "จ่ายแล้ว",
				Data:    nil,
			})
		case "NOT_READY":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "session ไม่พร้อมคิดเงิน",
				Data:    nil,
			})
		case "INSUFFICIENT":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "receivedAmount ไม่พอ",
				Data:    nil,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
				Success: false,
				Message: "เกิดข้อผิดพลาดภายในระบบ",
				Data:    nil,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "ชำระเงินสำเร็จ",
		Data:    resp,
	})
}

func (h *PaymentHandler) GetByID(c *fiber.Ctx) error {
	paymentID, err := strconv.Atoi(c.Params("paymentId"))
	if err != nil || paymentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.GetByID(c.UserContext(), paymentID)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "ข้อมูลไม่ถูกต้อง",
				Data:    nil,
			})
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ payment",
				Data:    nil,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
				Success: false,
				Message: "เกิดข้อผิดพลาดภายในระบบ",
				Data:    nil,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงข้อมูลการชำระเงินสำเร็จ",
		Data:    resp,
	})
}

func (h *PaymentHandler) GetAll(c *fiber.Ctx) error {
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	paymentMethodIDStr := c.Query("paymentMethodId")
	var paymentMethodID *int
	if paymentMethodIDStr != "" {
		id, err := strconv.Atoi(paymentMethodIDStr)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "ข้อมูลไม่ถูกต้อง",
				Data:    nil,
			})
		}
		paymentMethodID = &id
	}

	filter := models.PaymentListFilter{
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		PaymentMethodID: paymentMethodID,
		Status:          status,
		Page:            page,
		Limit:           limit,
	}

	items, total, err := h.service.GetAll(c.UserContext(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายการการชำระเงินสำเร็จ",
		Data: fiber.Map{
			"items": items,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}
