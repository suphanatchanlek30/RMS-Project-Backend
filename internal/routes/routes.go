package routes

import (
	"github.com/suphanatchanlek30/rms-project-backend/internal/handlers"
	"github.com/suphanatchanlek30/rms-project-backend/internal/middleware"
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

	roleRepo := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := handlers.NewRoleHandler(roleService)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", middleware.Protected(), authHandler.Me)
	auth.Post("/logout", middleware.Protected(), authHandler.Logout)

	v1.Get("/tables", tableHandler.GetAll)
	v1.Get("/customer/menus", menuHandler.GetCustomerMenus)
	v1.Get("/roles", middleware.Protected(), middleware.AdminOnly(), roleHandler.GetAll)

	v1.Post("/employees", middleware.Protected(), middleware.AdminOnly(), employeeHandler.CreateEmployee)
	v1.Get("/employees", middleware.Protected(), middleware.AdminOnly(), employeeHandler.GetEmployees)
}
