package models

type Product struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	Name  string  `json:"name" gorm:"not null"`
	Price float64 `json:"price" gorm:"not null"`
}
