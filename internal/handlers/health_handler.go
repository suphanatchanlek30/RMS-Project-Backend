package handlers

import (
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "server is running",
		Data: map[string]string{
			"status": "ok",
		},
	})
}
