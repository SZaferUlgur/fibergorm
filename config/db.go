package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	mysqlerr "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB paket genelinde kullanılacak global GORM veritabanı bağlantısı
var DB *gorm.DB

func ConnectDB() {
	// .env dosyasından gereklileri çekelim
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// DSN (Data Source Name) veritabanı baglantısı
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	var err error
	// GORM kullanarak veritabanına baglan
	level := os.Getenv("GORM_LOG_LEVEL")

	var logLevel logger.LogLevel

	switch strings.ToLower(level) {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	default:
		logLevel = logger.Silent
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	// Bağlantı hatası varsa özel olarak veritabanı bulunmaması hatası 1049
	if err != nil {
		mysqlErr, ok := err.(*mysqlerr.MySQLError)
		if ok && mysqlErr.Number == 1049 {
			log.Printf("Veritabanı '%s' bulunamadı. Oluşturuluyor...", dbName)

			// Veritabanı yoksa oluştur
			if err := createDatabase(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
				log.Fatalf("Veritabanı oluşturulamadı: %v", err)
			}

			// Veritabanı oluşturulduktan sonra tekrar bağlanmayı dene
			DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				log.Fatalf("Yeni oluşturulan veritabanına bağlanılamadı: %v", err)
			}

			log.Println("Veritabanı başarıyla oluşturuldu ve bağlantı kuruldu.")

			// Oluşturulan veritabanında tabloları ve diğer yapıları oluşturmak için SQL scriptini çalıştır
			if err := runSqlScript(DB); err != nil {
				log.Fatalf("SQL script çalıştırılamadı: %v", err)
			}
			log.Println("Tablolar başarıyla oluşturuldu.")
		} else {
			// Veritabanı dışındaki farklı bağlantı hatalarında programı sonlandır
			log.Fatalf("Veritabanına bağlanırken kritik bir hata oluştu: %v", err)
		}
	} else {
		// Bağlantı başarılı ise bunu bildir
		log.Println("Veritabanı bağlantısı başarılı.")
	}
}

// create database fonksiyonu
func createDatabase(user, pass, host, port, dbName string) error {
	// veritabanı adı olmadan bağlantı oluştur (yeni ekleneceği için)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port,
	)

	// GORM ile bağlantı aç
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err // Bağlantı açılmazsa hata döndür
	}

	// create database sorgusu
	query := fmt.Sprintf(
		"CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
		dbName,
	)
	if err := db.Exec(query).Error; err != nil {
		return err // Oluşturma hatası varsa geri döndür
	}

	// bağlantıyı kapat / kaynakları serbest bırakmak için
	sqlDB, _ := db.DB()
	sqlDB.Close()

	return nil
}

// runSqlScript fonksiyonu
func runSqlScript(db *gorm.DB) error {
	// dosyanın yolu
	filePath := "database/script.sql"

	// dosyayı oku
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("SQL script dosyası okunamadı '%s': %w", filePath, err)
	}

	// sql dosyasını özel bir ayraç ile parçala..
	parts := strings.Split(string(content), "--||--")

	// for next döngüsü
	for i, query := range parts {
		query = strings.TrimSpace(query)
		if query == "" {
			continue // Boş sorguları atla
		}

		// bazı özel satırları da atla
		if strings.HasPrefix(strings.ToUpper(query), "DELIMITER") ||
			strings.HasPrefix(query, "/*!") ||
			strings.HasPrefix(strings.ToUpper(query), "USE ") {
			continue
		}

		// sorguyu loglama yapalım
		log.Printf("Running query part %d:\n%s\n", i+1, query)

		// sorguyu çalıştıralım.
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("SQL query part %d çalıştırılırken hata: %w\nQuery: %s", i+1, err, query)
		}
	}
	return nil // tüm sorgular başarılıysa nil döndür

}
