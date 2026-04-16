package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type QRSessionHandler struct {
	service *services.QRSessionService
}

func NewQRSessionHandler(service *services.QRSessionService) *QRSessionHandler {
	return &QRSessionHandler{service: service}
}

func (h *QRSessionHandler) CreateQRSession(c *fiber.Ctx) error {
	var req models.CreateQRSessionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.SessionID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.CreateQRSession(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ session",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "มี QR active อยู่แล้ว",
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
		Message: "สร้าง QR Session สำเร็จ",
		Data:    resp,
	})
}
