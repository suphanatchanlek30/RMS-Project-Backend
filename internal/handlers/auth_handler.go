package handlers

import (
	"errors"
	"strings"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "กรุณากรอก email และ password ให้ครบ",
			Data:    nil,
		})
	}

	resp, err := h.service.Login(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Message: "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
				Data:    nil,
			})
		}

		if errors.Is(err, repositories.ErrAccountDisabled) {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "บัญชีถูกปิดใช้งาน",
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
		Message: "เข้าสู่ระบบสำเร็จ",
		Data:    resp,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	employeeIDVal := c.Locals("employeeId")
	employeeID, ok := employeeIDVal.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Message: "token ไม่ถูกต้องหรือหมดอายุ",
			Data:    nil,
		})
	}

	resp, err := h.service.GetMe(c.UserContext(), employeeID)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidToken) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Message: "token ไม่ถูกต้องหรือหมดอายุ",
				Data:    nil,
			})
		}

		if errors.Is(err, repositories.ErrAccountDisabled) {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "บัญชีถูกปิดใช้งาน",
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
		Message: "ดึงข้อมูลผู้ใช้สำเร็จ",
		Data:    resp,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	err := h.service.Logout(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "ออกจากระบบไม่สำเร็จ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ออกจากระบบสำเร็จ",
		Data:    nil,
	})
}
