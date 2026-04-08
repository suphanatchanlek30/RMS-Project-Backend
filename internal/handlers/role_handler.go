package handlers

import (
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(service *services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetAll(c *fiber.Ctx) error {
	roles, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายการ role สำเร็จ",
		Data:    roles,
	})
}