package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type TableSessionHandler struct {
	service *services.TableSessionService
}

func NewTableSessionHandler(service *services.TableSessionService) *TableSessionHandler {
	return &TableSessionHandler{service: service}
}

func (h *TableSessionHandler) OpenTable(c *fiber.Ctx) error {
	var req models.OpenTableRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.TableID <= 0 || req.EmployeeID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.OpenTable(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบโต๊ะ",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "โต๊ะกำลังใช้งานอยู่",
				Data:    nil,
			})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "เปิดโต๊ะไม่ได้ตามกฎธุรกิจ",
				Data:    nil,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
				Success: false,
				Message: "เปิดโต๊ะไม่สำเร็จ",
				Data:    nil,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "เปิดโต๊ะสำเร็จ",
		Data:    resp,
	})
}
