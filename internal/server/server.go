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

	// dsn := "postgres://username:password@db.supabase_url:5432/dbname?sslmode=require"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	PrepareStmt: false, // Nonaktifkan prepared statement cache
	// })
	// if err != nil {
	// 	panic("failed to connect database")
	// }

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false, // Nonaktifkan prepared statement cache
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// err = database.Migrator().DropTable(&Post{}, &User{})
	// if err != nil {
	// 	log.Fatalf("failed to drop table: %v", err)
	// }

	// err = database.AutoMigrate(&Post{}, &User{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to auto migrate: %w", err)
	// }

	log.Println("🚀 Successfully connected to database & migration completed!")

	return database, nil
}
