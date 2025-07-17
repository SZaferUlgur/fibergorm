package routes

import (
	"errors"
	"fibergorm/config"
	"fibergorm/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	app.Post("/api/v1/login", func(c *fiber.Ctx) error {
		return loginHandler(c, db)
	})
}

// login fonksiyonu
func loginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz İstek Gövdesi",
		})
	}

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Geçersiz EPosta ve Şifre",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı Bulunamadı",
		})
	}

	// şifre karşılaştırması
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Geçersiz EPosta ve Şifre",
		})
	}

	// jwt token oluştur
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token Oluşturma Hatası",
		})
	}

	// herşey tamam ise
	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
