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
	resp, err := h.service.GetTablesOverview(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงภาพรวมโต๊ะสำเร็จ",
		Data:    resp,
	})
}

func (h *CashierHandler) GetSessionCheckout(c *fiber.Ctx) error {
	sessionID, err := c.ParamsInt("sessionId")
	if err != nil || sessionID <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
			Success: false,
			Message: "sessionId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.GetSessionCheckout(c.UserContext(), sessionID)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ session",
				Data:    nil,
			})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "ข้อมูลไม่ถูกต้อง",
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
		Message: "ดึงข้อมูล checkout สำเร็จ",
		Data:    resp,
	})
}

func (h *CashierHandler) Checkout(c *fiber.Ctx) error {
	var req models.CashierCheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.Checkout(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ session หรือ payment method",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "session ถูกปิดไปแล้วหรือชำระเงินแล้ว",
				Data:    nil,
			})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "ข้อมูลไม่ถูกต้อง",
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
