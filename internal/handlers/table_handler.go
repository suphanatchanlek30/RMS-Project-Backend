package handlers

import (
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type TableHandler struct {
	service *services.TableService
}

func NewTableHandler(service *services.TableService) *TableHandler {
	return &TableHandler{service: service}
}

func (h *TableHandler) GetAll(c *fiber.Ctx) error {
	tables, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "failed to fetch tables",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "fetch tables success",
		Data:    tables,
	})
}
