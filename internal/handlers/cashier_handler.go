package handlers

import (
	"rms-project-backend/internal/models"
	"rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CashierHandler struct {
	cashierService services.CashierService
}

func NewCashierHandler(cashierService services.CashierService) *CashierHandler {
	return &CashierHandler{
		cashierService: cashierService,
	}
}

func (h *CashierHandler) GetTablesOverview(c *fiber.Ctx) error {
	data, err := h.cashierService.GetTablesOverview(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ดึงภาพรวมโต๊ะไม่สำเร็จ",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.CashierTablesOverviewResponse{
		Success: true,
		Message: "ดึงภาพรวมโต๊ะสำเร็จ",
		Data:    data,
	})
}