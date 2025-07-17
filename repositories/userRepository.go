package repositories

import (
	"fibergorm/models"

	"gorm.io/gorm"
)

// userrepository interface
type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

// gorm bağlantısı ile işlemler yapılacak.
// userrepo'nun somut struct'ı
type userRepository struct {
	db *gorm.DB
}

// yeni bir userrepo oluşturmak
// db parametresi gorm veritabanı bağlantısı
// bağımlılık injection sağlanacak
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User     // boş bir user dizisi
	result := r.db.Find(&users) // gorm ile tümünü çek
	return users, result.Error  // kullanıcı listesi veya hata dönder

}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id) // gorm ile id'ye göre ilk eşleşme verilerini çek
	return &user, result.Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user) // parametre olarak email'i soru isaretiyle arat
	// select * from users where email = p_email(?) limit 1;
	return &user, result.Error
}

func (r *userRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *userRepository) UpdateUser(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *userRepository) DeleteUser(id uint) error {
	result := r.db.Delete(&models.User{}, id)

	// eğer silinecek kayıt yoksa rows_effected == 0 ROW_COUNT() > 0
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
