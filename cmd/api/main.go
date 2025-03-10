package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/routes"
	"github.com/rianabd01/socialblog-be/internal/server"
)

func main() {
	r := gin.Default()

	db, err := server.ConnectDatabase()
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	// conn to redis
	if err := server.InitRedis(); err != nil {
		log.Println("Redis is not available, but the application will continue running.")
	}

	server.DB = db

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// routes
	routes.AuthRoutes(r)
	routes.BlogRoutes(r)
	routes.PostRoutes(r)
	routes.UserActionRoutes(r)

	r.Run(":8080")
}
