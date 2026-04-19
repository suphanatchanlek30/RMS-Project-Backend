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

	roleRepo := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := handlers.NewRoleHandler(roleService)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	tableSessionRepo := repositories.NewTableSessionRepository(db)
	tableSessionService := services.NewTableSessionService(tableSessionRepo)
	tableSessionHandler := handlers.NewTableSessionHandler(tableSessionService)

	qrSessionRepo := repositories.NewQRSessionRepository(db)
	qrSessionService := services.NewQRSessionService(qrSessionRepo, tableSessionRepo)
	qrSessionHandler := handlers.NewQRSessionHandler(qrSessionService)

	menuRepo := repositories.NewMenuRepository(db)
	menuService := services.NewMenuService(menuRepo, qrSessionRepo)
	menuHandler := handlers.NewMenuHandler(menuService)

	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo, qrSessionRepo, tableSessionRepo, tableRepo, menuRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	kitchenRepo := repositories.NewKitchenRepository(db)
	kitchenService := services.NewKitchenService(kitchenRepo)
	kitchenHandler := handlers.NewKitchenHandler(kitchenService)

	paymentRepo := repositories.NewPaymentRepository(db)
	receiptRepo := repositories.NewReceiptRepository(db)
	paymentService := services.NewPaymentService(paymentRepo, tableSessionRepo, receiptRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	receiptService := services.NewReceiptService(receiptRepo, paymentRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptService)

	paymentMethodRepo := repositories.NewPaymentMethodRepository(db)
	paymentMethodService := services.NewPaymentMethodService(paymentMethodRepo)
	paymentMethodHandler := handlers.NewPaymentMethodHandler(paymentMethodService)

	cashierRepo := repositories.NewCashierRepository(db)
	cashierService := services.NewCashierService(cashierRepo)
	cashierHandler := handlers.NewCashierHandler(cashierService)

	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", middleware.Protected(), authHandler.Me)
	auth.Post("/logout", middleware.Protected(), authHandler.Logout)

	v1.Get("/tables", middleware.Protected(), middleware.AdminOrCashier(), tableHandler.GetAll)
	v1.Get("/tables/:tableId", middleware.Protected(), middleware.AdminOrCashier(), tableHandler.GetByID)
	v1.Post("/tables", middleware.Protected(), middleware.AdminOnly(), tableHandler.Create)
	v1.Patch("/tables/:tableId", middleware.Protected(), middleware.AdminOnly(), tableHandler.Update)

	v1.Get("/customer/menus", menuHandler.GetCustomerMenus)
	v1.Get("/roles", middleware.Protected(), middleware.AdminOnly(), roleHandler.GetAll)

	v1.Post("/employees", middleware.Protected(), middleware.AdminOnly(), employeeHandler.CreateEmployee)
	v1.Get("/employees", middleware.Protected(), middleware.AdminOnly(), employeeHandler.GetEmployees)
	v1.Get("/employees/:employeeId", middleware.Protected(), middleware.AdminOnly(), employeeHandler.GetEmployeeByID)
	v1.Patch("/employees/:employeeId", middleware.Protected(), middleware.AdminOnly(), employeeHandler.UpdateEmployee)
	v1.Patch("/employees/:employeeId/status", middleware.Protected(), middleware.AdminOnly(), employeeHandler.UpdateEmployeeStatus)

	v1.Post("/table-sessions/open", middleware.Protected(), middleware.CashierOnly(), tableSessionHandler.OpenTable)
	v1.Get("/table-sessions/:sessionId", middleware.Protected(), middleware.AdminOrCashier(), tableSessionHandler.GetByID)
	v1.Get("/table-sessions/:sessionId/bill", middleware.Protected(), middleware.AdminOrCashier(), tableSessionHandler.GetSessionBill)
	v1.Get("/tables/:tableId/current-session", middleware.Protected(), middleware.AdminOrCashier(), tableSessionHandler.GetCurrentByTableID)
	v1.Patch("/table-sessions/:sessionId/close", middleware.Protected(), middleware.CashierOnly(), tableSessionHandler.CloseSession)

	v1.Post("/qr-sessions", middleware.Protected(), middleware.CashierOnly(), qrSessionHandler.CreateQRSession)
	v1.Get("/qr-sessions/:qrSessionId", middleware.Protected(), middleware.AdminOrCashier(), qrSessionHandler.GetByID)
	v1.Get("/qr/:token", qrSessionHandler.VerifyQR)

	v1.Post("/customer/orders", orderHandler.CreateCustomerOrder)
	v1.Get("/customer/orders", orderHandler.GetCustomerOrders)
	v1.Post("/orders", middleware.Protected(), middleware.CashierOnly(), orderHandler.CreateOrder)
	v1.Get("/orders/:orderId", middleware.Protected(), middleware.AdminCashierChef(), orderHandler.GetByID)
	v1.Get("/table-sessions/:sessionId/orders", middleware.Protected(), middleware.AdminOrCashier(), orderHandler.GetBySessionID)

	v1.Post("/categories", middleware.Protected(), middleware.AdminOnly(), categoryHandler.Create)
	v1.Get("/categories", categoryHandler.GetAll)
	v1.Patch("/categories/:categoryId", middleware.Protected(), middleware.AdminOnly(), categoryHandler.Update)

	v1.Get("/menus", middleware.Protected(), middleware.AdminOrCashier(), menuHandler.GetAll)
	v1.Get("/menus/:menuId", middleware.Protected(), middleware.AdminOrCashier(), menuHandler.GetByID)
	v1.Post("/menus", middleware.Protected(), middleware.AdminOnly(), menuHandler.Create)
	v1.Patch("/menus/:menuId", middleware.Protected(), middleware.AdminOnly(), menuHandler.Update)
	v1.Patch("/menus/:menuId/status", middleware.Protected(), middleware.AdminOnly(), menuHandler.UpdateStatus)

	v1.Get("/orders/:orderId/items", middleware.Protected(), middleware.AdminCashierChef(), orderHandler.GetOrderItems)
	v1.Patch("/order-items/:orderItemId", middleware.Protected(), middleware.CashierOnly(), orderHandler.UpdateOrderItemQuantity)
	v1.Delete("/order-items/:orderItemId", middleware.Protected(), middleware.CashierOnly(), orderHandler.CancelOrderItem)

	v1.Get("/kitchen/orders", middleware.Protected(), middleware.ChefOnly(), kitchenHandler.GetKitchenOrders)
	v1.Patch("/order-items/:orderItemId/status", middleware.Protected(), middleware.ChefOnly(), orderHandler.UpdateOrderItemStatus)
	v1.Get("/order-items/:orderItemId/history", middleware.Protected(), middleware.AdminCashierChef(), orderHandler.GetOrderItemStatusHistory)
	v1.Get("/customer/order-status", orderHandler.GetCustomerOrderStatus)

	v1.Post("/payments", middleware.Protected(), middleware.CashierOnly(), paymentHandler.Create)
	v1.Get("/payments/:paymentId", middleware.Protected(), middleware.AdminOrCashier(), paymentHandler.GetByID)
	v1.Get("/payments", middleware.Protected(), middleware.AdminOnly(), paymentHandler.GetAll)
	v1.Get("/payment-methods", middleware.Protected(), middleware.AdminOrCashier(), paymentMethodHandler.GetAll)
	v1.Get("/payments/:paymentId/receipt", middleware.Protected(), middleware.AdminOrCashier(), receiptHandler.GetByPaymentID)
	v1.Get("/receipts/:receiptId", middleware.Protected(), middleware.AdminOrCashier(), receiptHandler.GetByReceiptID)

	v1.Get("/cashier/tables/overview", middleware.Protected(), middleware.CashierOnly(), cashierHandler.GetTablesOverview)
	v1.Get("/cashier/sessions/:sessionId/checkout", middleware.Protected(), middleware.CashierOnly(), cashierHandler.GetCheckout)
}
