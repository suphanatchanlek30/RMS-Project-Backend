package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(s *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req models.CreateEmployeeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.EmployeeName == "" || req.RoleID == 0 ||
		req.PhoneNumber == "" || req.Email == "" ||
		req.HireDate == "" || req.Password == "" {

		return c.Status(400).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ครบ",
			Data:    nil,
		})
	}

	resp, err := h.service.CreateEmployee(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "EMAIL_EXISTS":
			return c.Status(409).JSON(models.APIResponse{
				Success: false,
				Message: "email ซ้ำ",
				Data:    nil,
			})
		case "ROLE_NOT_FOUND":
			return c.Status(404).JSON(models.APIResponse{
				Success: false,
				Message: "role ไม่พบ",
				Data:    nil,
			})
		default:
			return c.Status(500).JSON(models.APIResponse{
				Success: false,
				Message: "เกิดข้อผิดพลาดภายในระบบ",
				Data:    nil,
			})
		}
	}

	return c.Status(201).JSON(models.APIResponse{
		Success: true,
		Message: "สร้างพนักงานสำเร็จ",
		Data:    resp,
	})
}

func (h *EmployeeHandler) GetEmployees(c *fiber.Ctx) error {

	roleIDStr := c.Query("roleId")
	statusStr := c.Query("status")
	search := c.Query("search")

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	var roleID *int
	if roleIDStr != "" {
		id, _ := strconv.Atoi(roleIDStr)
		roleID = &id
	}

	var status *bool
	if statusStr != "" {
		val := statusStr == "true"
		status = &val
	}

	items, total, err := h.service.GetEmployees(
		c.UserContext(),
		roleID,
		status,
		search,
		page,
		limit,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "server error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ดึงรายการพนักงานสำเร็จ",
		"data": fiber.Map{
			"items": items,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}
