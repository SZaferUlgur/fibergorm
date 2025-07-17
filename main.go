package main

import (
	"fibergorm/config"
	"fibergorm/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	config.ConnectDB()

	// yeni fiber uygulaması başlatıldı (HTTP SUNUCU)
	app := fiber.New()

	// Kullanıcı ile ilgili API Rotaları ve Diğerieri
	routes.SetupUserRoutes(app, config.DB)
	routes.SetupProductsRoutes(app, config.DB)
	routes.SetupAuthRoutes(app, config.DB)

	port := os.Getenv("FIBER_PORT")
	if port != "" {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		if err := app.Listen(":5000"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}

}
