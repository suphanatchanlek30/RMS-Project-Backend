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

func (h *OrderHandler) GetOrderItems(c *fiber.Ctx) error {
	orderIDParam := c.Params("orderId")

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "orderId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	items, err := h.service.GetOrderItems(c.Context(), orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ order",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาด",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายการอาหารสำเร็จ",
		Data:    items,
	})
}

func (h *OrderHandler) UpdateOrderItemQuantity(c *fiber.Ctx) error {
	idParam := c.Params("orderItemId")

	orderItemID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "orderItemId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	var req models.UpdateOrderItemRequest
	if err := c.BodyParser(&req); err != nil || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	result, err := h.service.UpdateOrderItemQuantity(c.Context(), orderItemID, req.Quantity)
	if err != nil {
		if err.Error() == "not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ order item",
				Data:    nil,
			})
		}
		if err.Error() == "invalid status" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "สถานะไม่อนุญาตให้แก้",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาด",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "แก้ไขจำนวนรายการอาหารสำเร็จ",
		Data:    result,
	})
}

func (h *OrderHandler) CancelOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("orderItemId")

	orderItemID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "orderItemId ไม่ถูกต้อง",
			Data:    nil,
		})
	}

	result, err := h.service.CancelOrderItem(c.UserContext(), orderItemID)
	if err != nil {

		if err.Error() == "order item not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบรายการอาหาร",
				Data:    nil,
			})
		}

		if err.Error() == "ไม่สามารถยกเลิกรายการนี้ได้" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Message: "รายการถูกทำแล้วหรือชำระแล้ว",
				Data:    nil,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "ยกเลิกรายการอาหารสำเร็จ",
		Data:    result,
	})
}

func (h *OrderHandler) UpdateOrderItemStatus(c *fiber.Ctx) error {
	orderItemID, _ := strconv.Atoi(c.Params("orderItemId"))

	var req models.UpdateOrderItemStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
		})
	}

	chefIDVal := c.Locals("employeeId")
	chefID, ok := chefIDVal.(int)
	if !ok {
		return c.Status(401).JSON(models.APIResponse{
			Success: false,
			Message: "token ไม่ถูกต้อง",
		})
	}

	resp, err := h.service.UpdateOrderItemStatus(
		c.UserContext(),
		orderItemID,
		req.Status,
		chefID,
	)

	if err != nil {
		return c.Status(422).JSON(models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "อัปเดตสถานะอาหารสำเร็จ",
		Data:    resp,
	})
}

func (h *OrderHandler) GetOrderItemStatusHistory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("orderItemId"))

	result, err := h.service.GetOrderItemStatusHistory(c.UserContext(), id)
	if err != nil {
		return c.Status(404).JSON(models.APIResponse{
			Success: false,
			Message: "ไม่พบข้อมูล",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "ดึงประวัติสถานะสำเร็จ",
		Data:    result,
	})
}

func (h *OrderHandler) GetCustomerOrderStatus(c *fiber.Ctx) error {
	qrToken := c.Query("qrToken")

	if qrToken == "" {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Message: "qrToken จำเป็น",
		})
	}

	resp, err := h.service.GetCustomerOrderStatus(c.UserContext(), qrToken)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(404).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบ QR หรือคำสั่งซื้อ",
				Data:    nil,
			})
		case "GONE":
			return c.Status(410).JSON(models.APIResponse{
				Success: false,
				Message: "QR หมดอายุ",
				Data:    nil,
			})
		case "UNPROCESSABLE":
			return c.Status(422).JSON(models.APIResponse{
				Success: false,
				Message: "session ปิดแล้ว",
				Data:    nil,
			})
		default:
			return c.Status(500).JSON(models.APIResponse{
				Success: false,
				Message: "เกิดข้อผิดพลาด",
				Data:    nil,
			})
		}
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "ดึงสถานะออเดอร์สำเร็จ",
		Data:    resp,
	})
}
