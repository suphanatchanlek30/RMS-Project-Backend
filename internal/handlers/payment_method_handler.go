package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type PaymentMethodHandler struct {
	service *services.PaymentMethodService
}

func NewPaymentMethodHandler(service *services.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{service: service}
}

func (h *PaymentMethodHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงวิธีชำระเงินสำเร็จ",
		Data:    items,
	})
}
