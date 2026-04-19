package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetSalesReport(c *fiber.Ctx) error {
	query, err := services.ParseSalesReportQuery(
		c.Query("dateFrom"),
		c.Query("dateTo"),
		c.Query("groupBy"),
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "query ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	items, err := h.service.GetSalesReport(c.UserContext(), query)
	if err != nil {
		if err.Error() == "BAD_REQUEST" {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "query ไม่ถูกต้อง",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายงานยอดขายสำเร็จ",
		Data:    items,
	})
}

func (h *ReportHandler) GetTopMenusReport(c *fiber.Ctx) error {
	limit := 10
	limitRaw := c.Query("limit")
	if limitRaw != "" {
		parsedLimit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "query ไม่ถูกต้อง",
				Data:    nil,
			})
		}
		limit = parsedLimit
	}

	query, err := services.ParseTopMenusReportQuery(
		c.Query("dateFrom"),
		c.Query("dateTo"),
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "query ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	items, err := h.service.GetTopMenusReport(c.UserContext(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายงานเมนูขายดีสำเร็จ",
		Data:    items,
	})
}
