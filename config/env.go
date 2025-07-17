package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string

func LoadEnv() {
	// ortam değişkenlerini alma
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".Env dosyası bulunamadı")
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET ortam değişkeni bulunamadı")
	}
}
