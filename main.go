package main

import (
	"log"

	"github.com/rianabd01/socialblog-be/controllers/postcontroller"
	"github.com/rianabd01/socialblog-be/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db, err := models.ConnectDatabase()
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	models.DB = db // ✅

	// posts
	r.GET("/api/posts", postcontroller.Index)
	r.POST("/api/posts", postcontroller.Create)
	r.GET("/api/posts/:id", postcontroller.ShowDetail)
	r.PUT("/api/posts/:id", postcontroller.Update)

	r.Run()
}
