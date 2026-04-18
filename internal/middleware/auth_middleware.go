package middleware

import (
	"strings"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่ได้ login หรือไม่มี token",
				Data:    nil,
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Message: "token ไม่ถูกต้อง",
				Data:    nil,
			})
		}

		claims, err := utils.ParseJWT(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Message: "token ไม่ถูกต้องหรือหมดอายุ",
				Data:    nil,
			})
		}

		c.Locals("employeeId", claims.EmployeeID)
		c.Locals("roleId", claims.RoleID)
		c.Locals("roleName", claims.RoleName)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, ok := c.Locals("roleName").(string)
		if !ok || roleName == "" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึงข้อมูล role",
				Data:    nil,
			})
		}

		if roleName != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึงข้อมูล role",
				Data:    nil,
			})
		}

		return c.Next()
	}
}

func AdminOrCashier() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, ok := c.Locals("roleName").(string)
		if !ok || roleName == "" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		if roleName != "ADMIN" && roleName != "CASHIER" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		return c.Next()
	}
}

func CashierOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, ok := c.Locals("roleName").(string)
		if !ok || roleName == "" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		if roleName != "CASHIER" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		return c.Next()
	}
}

func AdminCashierChef() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, ok := c.Locals("roleName").(string)
		if !ok || roleName == "" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		if roleName != "ADMIN" && roleName != "CASHIER" && roleName != "CHEF" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		return c.Next()
	}
}

func ChefOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, ok := c.Locals("roleName").(string)
		if !ok || roleName == "" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		if strings.ToUpper(roleName) != "CHEF" {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่มีสิทธิ์เข้าถึง",
				Data:    nil,
			})
		}

		return c.Next()
	}
}
