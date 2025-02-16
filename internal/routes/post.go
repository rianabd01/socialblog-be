package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/controllers/postcontroller"
)

func PostRoutes(r *gin.Engine) {
	r.GET("/api/posts", postcontroller.Index)
	r.POST("/api/posts", postcontroller.Create)
	r.GET("/api/posts/:id", postcontroller.ShowDetail)
	r.PUT("/api/posts/:id", postcontroller.Update)
}
