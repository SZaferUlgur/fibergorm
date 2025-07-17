package models

type User struct {
	// ID alanı birincil anahtar işaretlenir
	// uint tipi pozitif ve sayıları temsil eder
	ID uint `gorm:"primaryKey" json:"id"`
	// not null ifadesi veritabanına insert edilirken boş olamaz
	Name string `json:"name" gorm:"not null"`
	// unique ifadesi veritabanına insert edilirken benzersiz olmalı
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
}
