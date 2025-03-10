package server

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
	ctx         = context.Background() // Gunakan context global untuk Redis
)

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
	// err = db.Migrator().DropTable(&models.Blog{}, &models.User{}, &models.Post{}, &models.Like{})
	// if err != nil {
	// 	log.Fatalf("failed to drop table: %v", err)
	// }

	// err = db.AutoMigrate(&models.Blog{}, &models.User{}, &models.Post{}, &models.Like{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to auto migrate: %w", err)
	// }

	return db, nil
}

// InitRedis menghubungkan aplikasi dengan Redis
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Sesuaikan dengan alamat Redis kamu
		Password: "",               // Jika Redis memiliki password, masukkan di sini
		DB:       0,                // Gunakan DB default (0)
	})

	// Cek koneksi Redis
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Println("Warning: Failed to connect to Redis:", err)
		return err
	}

	log.Println("Connected to Redis successfully")
	return nil
}
