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

func (h *TableSessionHandler) GetByID(c *fiber.Ctx) error {
	sessionID, err := c.ParamsInt("sessionId")
	if err != nil || sessionID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "sessionId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	session, err := h.service.GetByID(c.UserContext(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Message: "ไม่พบ session",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงข้อมูล session สำเร็จ",
		Data:    session,
	})
}

func (h *TableSessionHandler) GetCurrentByTableID(c *fiber.Ctx) error {
	tableID, err := c.ParamsInt("tableId")
	if err != nil || tableID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "tableId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	session, err := h.service.GetCurrentSessionByTableID(c.UserContext(), tableID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Message: "ไม่มี session ที่เปิดอยู่",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึง session ปัจจุบันสำเร็จ",
		Data:    session,
	})
}

func (h *TableSessionHandler) CloseSession(c *fiber.Ctx) error {
	sessionID, err := c.ParamsInt("sessionId")
	if err != nil || sessionID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "sessionId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.CloseSession(c.UserContext(), sessionID)
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
				Message: "ยังมีบิลค้างชำระ",
				Data:    nil,
			})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "session ปิดไปแล้ว",
				Data:    nil,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
				Success: false,
				Message: "ปิดโต๊ะไม่สำเร็จ",
				Data:    nil,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ปิดโต๊ะสำเร็จ",
		Data:    resp,
	})
}
