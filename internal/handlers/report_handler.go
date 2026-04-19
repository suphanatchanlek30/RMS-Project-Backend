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
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	groupBy := c.Query("groupBy")

	// Validate required parameters
	if dateFrom == "" || dateTo == "" || groupBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "กรุณาระบุ dateFrom, dateTo, และ groupBy",
			Data:    nil,
		})
	}

	report, err := h.service.GetSalesReport(c.UserContext(), dateFrom, dateTo, groupBy)
	if err != nil {
		switch err.Error() {
		case "INVALID_DATE_FROM", "INVALID_DATE_TO", "INVALID_DATE_RANGE", "INVALID_GROUP_BY":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "query ไม่ถูกต้อง",
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

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายงานยอดขายสำเร็จ",
		Data:    report,
	})
}

func (h *ReportHandler) GetTopMenusReport(c *fiber.Ctx) error {
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	limitQuery := c.Query("limit")

	if dateFrom == "" || dateTo == "" || limitQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "กรุณาระบุ dateFrom, dateTo, และ limit",
			Data:    nil,
		})
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "query ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	report, err := h.service.GetTopMenusReport(c.UserContext(), dateFrom, dateTo, limit)
	if err != nil {
		switch err.Error() {
		case "INVALID_DATE_FROM", "INVALID_DATE_TO", "INVALID_DATE_RANGE", "INVALID_LIMIT":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "query ไม่ถูกต้อง",
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

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายงานเมนูขายดีสำเร็จ",
		Data:    report,
	})
}
