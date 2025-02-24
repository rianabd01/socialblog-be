package server

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("user")
	pass := os.Getenv("password")
	host := os.Getenv("host")
	dbname := os.Getenv("dbname")
	port := os.Getenv("port")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", user, pass, host, port, dbname)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disable prepared statements
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// // aktifkan untuk migrasi, resikonya data semua terhapus
	// err = db.Migrator().DropTable(&models.Blog{}, &models.User{})
	// if err != nil {
	// 	log.Fatalf("failed to drop table: %v", err)
	// }

	// err = db.AutoMigrate(&models.Blog{}, &models.User{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to auto migrate: %w", err)
	// }

	return db, nil
}
