package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	resp, err := h.service.GetSummary(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงข้อมูล dashboard สำเร็จ",
		Data:    resp,
	})
}
