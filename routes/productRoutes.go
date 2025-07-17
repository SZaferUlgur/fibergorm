package routes

import (
	"fibergorm/handlers"
	"fibergorm/middlewares"
	"fibergorm/repositories"

	"github.com/gofiber/fiber/v2" // fiber paketini kullanmak için
	"gorm.io/gorm"                // gorm paketini kullanmak için
)

// product rotaları oluşturuldu
func SetupProductsRoutes(app *fiber.App, db *gorm.DB) {
	productRepository := repositories.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepository)

	// rotaları grupla
	api := app.Group("/api/v1/products")

	// middleware kontrolü
	api.Use(middlewares.AuthMiddleware)

	// 5000:/api/v1/products/
	api.Post("/", productHandler.CreateProduct)
	// 5000:/api/v1/products/123
	api.Put("/:id", productHandler.UpdateProduct)
	api.Delete("/:id", productHandler.DeleteProduct)
	api.Get("/", productHandler.GetAllProducts)
	api.Get("/:id", productHandler.GetProductByID)
}
