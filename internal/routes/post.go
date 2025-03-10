package routes

import (
	"github.com/gin-gonic/gin"
	postcontroller "github.com/rianabd01/socialblog-be/internal/controllers/post-controller"
	"github.com/rianabd01/socialblog-be/internal/middleware"
)

func PostRoutes(r *gin.Engine) {
	r.GET("/api/posts", postcontroller.Index)
	r.GET("/api/posts/:id", postcontroller.ShowDetail)
	r.POST("/api/posts", middleware.AuthMiddleware(), postcontroller.Create)
	r.PUT("/api/posts/:id", middleware.AuthMiddleware(), postcontroller.Update)
}
