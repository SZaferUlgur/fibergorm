package handlers

import (
	"errors"
	"fibergorm/models"
	"fibergorm/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userhandler için struct
type UserHandler struct {
	repo repositories.UserRepository
}

// newuser fonksiyonu
func NewUserHandler(repo repositories.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// get all users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcılar alınırken Hata Oluştu",
		})
	}

	// frontend'e gönderme
	return c.JSON(users)
}

// get user by id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz kullanıcı ID",
		})
	}

	user, err := h.repo.FindByID(uint(id))
	if err != nil {
		// eğer kullanıcı yoksa
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// frontend'e gönderme
	return c.JSON(user)

}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// model oluştur
	user := new(models.User)

	// istek gövdesini al
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz İstek Gövdesi",
		})
	}

	// şifre zorunlu ve boş olmasın
	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Sifre Zorunlu",
		})
	}

	// şifreyi hash edelim
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Şifreleme Hatası",
		})
	}

	// hash'lenmiş şifresi modele atama yapalım
	user.Password = string(hashPassword)

	// log.Default().Println(user)

	// veritabanına ekleme
	if err := h.repo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı Ekleme Hatası",
		})
	}

	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(user)
}

// update user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz kullanıcı ID",
		})
	}

	// önce kullanıdı db'de kontrol edelim
	user, err := h.repo.FindByID(uint(id))
	if err != nil {
		// 404
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Kullanıcı bulunamadı",
		})
	}

	// güncellenmiş verileri tutmak için yeni bir user oluştur
	updateUser := new(models.User)
	if err := c.BodyParser(updateUser); err != nil {
		// 400 döndürmek
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz İstek Gövdesi",
		})
	}

	// var olan kullanı alanlarını güncelle
	user.Name = updateUser.Name
	user.Email = updateUser.Email

	// repo üzerinden güncelle
	if err := h.repo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı Güncelleme Hatası",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// delete user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz kullanıcı ID",
		})
	}

	// repo üzerinden silme işlemi
	if err := h.repo.DeleteUser(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Kullanıcı bulunamadı",
			})
		}
		// başka hata 500 dönder
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı Silme Hatası",
		})
	}
	// başarılı ise bilgi dönder
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Kullanıcı Silindi",
	})
}
