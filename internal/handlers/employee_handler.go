package handlers

import (
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
