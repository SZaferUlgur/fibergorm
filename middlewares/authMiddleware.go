package middlewares

import (
	"fibergorm/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// header'dan token'i alalım
	tokenString := c.Get("token")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token bulunamadı",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Geçersiz token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user_id", claims["user_id"])
	}

	return c.Next()
}
