package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateCustomerOrder(c *fiber.Ctx) error {
	var req models.CreateCustomerOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	if req.QRToken == "" || len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	for _, item := range req.Items {
		if item.MenuID <= 0 || item.Quantity <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
		}
	}

	resp, err := h.service.CreateCustomerOrder(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบ QR หรือเมนู", Data: nil})
		case "GONE":
			return c.Status(fiber.StatusGone).JSON(models.APIResponse{Success: false, Message: "QR หมดอายุ", Data: nil})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{Success: false, Message: "เมนูปิดขาย/โต๊ะไม่พร้อมใช้งาน", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{Success: true, Message: "สร้างคำสั่งซื้อสำเร็จ", Data: resp})
}

func (h *OrderHandler) GetCustomerOrders(c *fiber.Ctx) error {
	qrToken := c.Query("qrToken")
	if qrToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "กรุณาระบุ qrToken", Data: nil})
	}

	resp, err := h.service.GetCustomerOrders(c.UserContext(), qrToken)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "กรุณาระบุ qrToken", Data: nil})
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบ QR", Data: nil})
		case "GONE":
			return c.Status(fiber.StatusGone).JSON(models.APIResponse{Success: false, Message: "QR หมดอายุ", Data: nil})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{Success: false, Message: "session ปิดแล้ว", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "ดึงคำสั่งซื้อของลูกค้าสำเร็จ", Data: resp})
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	if req.SessionID <= 0 || req.TableID <= 0 || req.CreatedByEmployeeID == nil || *req.CreatedByEmployeeID <= 0 || len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	for _, item := range req.Items {
		if item.MenuID <= 0 || item.Quantity <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
		}
	}

	resp, err := h.service.CreateOrder(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "BAD_REQUEST":
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบข้อมูลที่ต้องการ", Data: nil})
		case "UNPROCESSABLE":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{Success: false, Message: "เมนูปิดขาย/โต๊ะไม่พร้อมใช้งาน", Data: nil})
		case "INTERNAL":
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{Success: true, Message: "สร้างคำสั่งซื้อสำเร็จ", Data: resp})
}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil || orderID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	resp, err := h.service.GetOrderByID(c.UserContext(), orderID)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบ order", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "ดึงข้อมูล order สำเร็จ", Data: resp})
}

func (h *OrderHandler) GetBySessionID(c *fiber.Ctx) error {
	sessionID, err := strconv.Atoi(c.Params("sessionId"))
	if err != nil || sessionID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{Success: false, Message: "ข้อมูลไม่ถูกต้อง", Data: nil})
	}

	resp, err := h.service.GetOrdersBySessionID(c.UserContext(), sessionID)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{Success: false, Message: "ไม่พบ session", Data: nil})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{Success: false, Message: "เกิดข้อผิดพลาดภายในระบบ", Data: nil})
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{Success: true, Message: "ดึงรายการ order ของโต๊ะสำเร็จ", Data: resp})
}
