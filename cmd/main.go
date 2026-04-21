package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suphanatchanlek30/rms-project-backend/internal/config"
	"github.com/suphanatchanlek30/rms-project-backend/internal/database"
	"github.com/suphanatchanlek30/rms-project-backend/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.LoadEnv()
	appPort := config.GetEnv("APP_PORT", "8080")

	dbPool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer dbPool.Close()

	app := fiber.New(fiber.Config{
		AppName:      config.GetEnv("APP_NAME", "RMS Backend"),
		BodyLimit:    10 * 1024 * 1024,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		AllowMethods:     config.GetEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		AllowHeaders:     config.GetEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		AllowCredentials: config.GetEnv("CORS_ALLOW_CREDENTIALS", "false") == "true",
	}))

	routes.SetupRoutes(app, dbPool)

	go func() {
		if err := app.Listen(":" + appPort); err != nil {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	log.Printf("server is running on port %s", appPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}
