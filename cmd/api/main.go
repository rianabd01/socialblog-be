package main

import (
	"log"

	"github.com/rianabd01/socialblog-be/internal/routes"
	"github.com/rianabd01/socialblog-be/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db, err := server.ConnectDatabase()
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	server.DB = db

	// routes
	routes.AuthRoutes(r)
	routes.BlogRoutes(r)

	r.Run(":8080")
}
