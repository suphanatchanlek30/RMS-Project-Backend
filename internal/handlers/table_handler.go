package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type TableHandler struct {
	service *services.TableService
}

func NewTableHandler(service *services.TableService) *TableHandler {
	return &TableHandler{service: service}
}

func (h *TableHandler) GetAll(c *fiber.Ctx) error {
	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	tables, err := h.service.GetAll(c.UserContext(), status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "ดึงรายการโต๊ะไม่สำเร็จ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายการโต๊ะสำเร็จ",
		Data:    tables,
	})
}

func (h *TableHandler) GetByID(c *fiber.Ctx) error {
	tableID, err := c.ParamsInt("tableId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "tableId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	table, err := h.service.GetByID(c.UserContext(), tableID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Message: "ไม่พบโต๊ะ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงข้อมูลโต๊ะสำเร็จ",
		Data:    table,
	})
}

func (h *TableHandler) Create(c *fiber.Ctx) error {
	var req models.CreateTableRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.TableNumber == "" || req.Capacity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	table, err := h.service.Create(c.UserContext(), req.TableNumber, req.Capacity)
	if err != nil {

		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "หมายเลขโต๊ะซ้ำ",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "สร้างโต๊ะไม่สำเร็จ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "สร้างโต๊ะสำเร็จ",
		Data:    table,
	})
}
