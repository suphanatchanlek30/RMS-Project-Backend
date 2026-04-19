package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type CashierHandler struct {
	service *services.CashierService
}

func NewCashierHandler(service *services.CashierService) *CashierHandler {
	return &CashierHandler{service: service}
}

func (h *CashierHandler) GetTablesOverview(c *fiber.Ctx) error {
	tables, err := h.service.GetTablesOverview(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "ดึงภาพรวมโต๊ะไม่สำเร็จ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงภาพรวมโต๊ะสำเร็จ",
		Data:    tables,
	})
}

func (h *CashierHandler) GetCheckout(c *fiber.Ctx) error {
	sessionID, err := c.ParamsInt("sessionId")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
			Success: false,
			Message: "sessionId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	checkout, err := h.service.GetCheckout(c.UserContext(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Message: "ไม่พบข้อมูล checkout",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงข้อมูล checkout สำเร็จ",
		Data:    checkout,
	})
}

func (h *CashierHandler) Checkout(c *fiber.Ctx) error {
	var req models.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.Checkout(c.UserContext(), &req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ session",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "session ไม่พร้อมสำหรับการชำระเงิน",
				Data:    nil,
			})
		case "VALIDATION":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "จำนวนเงินที่รับไม่เพียงพอ",
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
		Message: "ชำระเงินและปิดโต๊ะสำเร็จ",
		Data:    resp,
	})
}