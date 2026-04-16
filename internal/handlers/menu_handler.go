package handlers

import (
	"strconv"

	"github.com/suphanatchanlek30/rms-project-backend/internal/models"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	service *services.MenuService
}

func NewMenuHandler(service *services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) GetCustomerMenus(c *fiber.Ctx) error {
	menus, err := h.service.GetCustomerMenus(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "failed to fetch customer menus",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "fetch customer menus success",
		Data:    menus,
	})
}

func (h *MenuHandler) GetAll(c *fiber.Ctx) error {
	categoryIDStr := c.Query("categoryId")
	keyword := c.Query("keyword")
	statusStr := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	var categoryID *int
	if categoryIDStr != "" {
		id, _ := strconv.Atoi(categoryIDStr)
		categoryID = &id
	}

	var status *bool
	if statusStr != "" {
		val := statusStr == "true"
		status = &val
	}

	items, total, err := h.service.GetAll(c.UserContext(), categoryID, keyword, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดภายในระบบ",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "ดึงรายการเมนูสำเร็จ",
		Data: fiber.Map{
			"items": items,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

func (h *MenuHandler) GetByID(c *fiber.Ctx) error {
	menuID, err := strconv.Atoi(c.Params("menuId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	m, err := h.service.GetByID(c.UserContext(), menuID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบเมนู",
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
		Message: "ดึงข้อมูลเมนูสำเร็จ",
		Data:    m,
	})
}

func (h *MenuHandler) Create(c *fiber.Ctx) error {
	var req models.CreateMenuRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.MenuName == "" || req.CategoryID <= 0 || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.Create(c.UserContext(), req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบหมวดหมู่",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "ชื่อเมนูซ้ำ",
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

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "สร้างเมนูสำเร็จ",
		Data:    resp,
	})
}

func (h *MenuHandler) Update(c *fiber.Ctx) error {
	menuID, err := strconv.Atoi(c.Params("menuId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	var req models.UpdateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	if req.MenuName == "" || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง",
			Data:    nil,
		})
	}

	resp, err := h.service.Update(c.UserContext(), menuID, req)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
				Success: false,
				Message: "ไม่พบเมนู",
				Data:    nil,
			})
		case "CONFLICT":
			return c.Status(fiber.StatusConflict).JSON(models.APIResponse{
				Success: false,
				Message: "ชื่อเมนูซ้ำ",
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
		Message: "อัปเดตเมนูสำเร็จ",
		Data:    resp,
	})
}
