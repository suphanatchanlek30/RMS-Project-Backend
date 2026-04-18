package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type ReceiptHandler struct {
	service *services.ReceiptService
}

func NewReceiptHandler(service *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{service: service}
}

func (h *ReceiptHandler) GetByPaymentID(c *fiber.Ctx) error {
	paymentID, err := strconv.Atoi(c.Params("paymentId"))
	if err != nil || paymentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	resp, err := h.service.GetByPaymentID(c.UserContext(), paymentID)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบ receipt หรือ payment", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "ดึงข้อมูลใบเสร็จสำเร็จ", Data: resp})
}

func (h *ReceiptHandler) GetByReceiptID(c *fiber.Ctx) error {
	receiptID, err := strconv.Atoi(c.Params("receiptId"))
	if err != nil || receiptID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	resp, err := h.service.GetByReceiptID(c.UserContext(), receiptID)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบ receipt", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "ดึงข้อมูลใบเสร็จสำเร็จ", Data: resp})
}
