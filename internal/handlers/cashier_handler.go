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