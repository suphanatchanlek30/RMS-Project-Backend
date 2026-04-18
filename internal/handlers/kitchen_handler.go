package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type KitchenHandler struct {
	service *services.KitchenService
}

func NewKitchenHandler(s *services.KitchenService) *KitchenHandler {
	return &KitchenHandler{service: s}
}

func (h *KitchenHandler) GetKitchenOrders(c *fiber.Ctx) error {
	status := c.Query("status")
	tableID, _ := strconv.Atoi(c.Query("tableId", "0"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.service.GetKitchenOrders(c.Context(), status, tableID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "ดึงคิวครัวสำเร็จ",
		Data:    result,
	})
}
