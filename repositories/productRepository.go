package repositories

import (
	"fibergorm/models"

	"gorm.io/gorm"
)

// interface
type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
}

// product repo interface uygulanır ve gorm ile islemler yapılır
type productRepository struct {
	db *gorm.DB
}

// yeni bir productrepo olusturmak
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	result := r.db.Find(&products)
	return products, result.Error
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	result := r.db.First(&product, id)
	return &product, result.Error
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	result := r.db.Create(product)
	return result.Error
}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	result := r.db.Save(product)
	return result.Error
}

func (r *productRepository) DeleteProduct(id uint) error {
	result := r.db.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
