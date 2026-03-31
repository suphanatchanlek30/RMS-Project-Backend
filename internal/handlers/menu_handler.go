package handlers

import (
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	service *services.MenuService
}

func NewMenuHandler(service *services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) GetCustomerMenus(c *fiber.Ctx) error {
	menus, err := h.service.GetCustomerMenus(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "failed to fetch customer menus",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "fetch customer menus success",
		Data:    menus,
	})
}
