package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
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

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Matikan prepared statement cache
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	sqlDB.SetConnMaxLifetime(0)
	sqlDB.Exec("DISCARD ALL")

	if os.Getenv("GIN_MODE") != "release" {

		if err := database.AutoMigrate(&Post{}); err != nil {
			log.Println("Skipping migration, table might already exist:", err)
		}
	}

	fmt.Println("Successfully connected to database!")

	DB = database
}
