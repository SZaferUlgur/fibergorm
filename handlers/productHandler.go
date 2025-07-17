package handlers

import (
	"errors"
	"fibergorm/models"
	"fibergorm/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type ProductHandler struct {
	repo repositories.ProductRepository
}

func NewProductHandler(repo repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Urunler alinirken hata olustu",
		})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz urun ID",
		})
	}
	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Urun bulunamadi",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(product)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz İstek Gövdesi",
		})
	}

	if err := h.repo.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Urun olusturulamadi",
		})
	}
	return c.Status(fiber.StatusOK).JSON(product)

}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz urun ID",
		})
	}

	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Urun bulunamadi",
		})
	}

	updateProduct := new(models.Product)
	if err := c.BodyParser(updateProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz İstek Gövdesi",
		})
	}

	product.Name = updateProduct.Name
	product.Price = updateProduct.Price

	// repo ile güncelle
	if err := h.repo.UpdateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Urun guncellenemedi",
		})
	}
	// fronend'e gönder
	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gecersiz urun ID",
		})
	}
	// repo ile sil
	if err := h.repo.DeleteProduct(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Urun bulunamadi",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Urun silindi",
	})
}
