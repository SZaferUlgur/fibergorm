package routes

import (
	"fibergorm/handlers"     // handler klasörü içindeki fonksiyonları kullanmak için
	"fibergorm/repositories" // repository klasörü içindeki fonksiyonları kullanmak için

	"github.com/gofiber/fiber/v2" // fiber paketini kullanmak için
	"gorm.io/gorm"                // gorm paketini kullanmak için
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
	// repository klasörü içindeki fonksiyonları kullanmak için
	userRepository := repositories.NewUserRepository(db)
	// handler klasörü içindeki fonksiyonları kullanmak için
	userHandler := handlers.NewUserHandler(userRepository)

	// Kullanıcı ile ilgili API Rotaları
	api := app.Group("/api/v1/users")
	// Rotalar
	api.Post("/", userHandler.CreateUser)      // Changed from fiber.Ctx to fiber.v2.Ctx
	api.Put("/:id", userHandler.UpdateUser)    // Changed from fiber.Ctx to fiber.v2.Ctx
	api.Delete("/:id", userHandler.DeleteUser) // Changed from fiber.Ctx to fiber.v2.Ctx
	api.Get("/", userHandler.GetAllUsers)      // Changed from fiber.Ctx to fiber.v2.Ctx
	api.Get("/:id", userHandler.GetUserByID)   // Changed from fiber.Ctx to fiber.v2.Ctx
}
