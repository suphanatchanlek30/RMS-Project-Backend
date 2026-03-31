package routes

import (
	"github.com/suphanatchanlek30/rms-project-backend/internal/handlers"
	"github.com/suphanatchanlek30/rms-project-backend/internal/repositories"
	"github.com/suphanatchanlek30/rms-project-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	healthHandler := handlers.NewHealthHandler()

	tableRepo := repositories.NewTableRepository(db)
	tableService := services.NewTableService(tableRepo)
	tableHandler := handlers.NewTableHandler(tableService)

	menuRepo := repositories.NewMenuRepository(db)
	menuService := services.NewMenuService(menuRepo)
	menuHandler := handlers.NewMenuHandler(menuService)

	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/tables", tableHandler.GetAll)
	v1.Get("/customer/menus", menuHandler.GetCustomerMenus)
}
