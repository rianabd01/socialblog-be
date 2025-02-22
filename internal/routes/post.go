package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/controllers/postcontroller"
	"github.com/rianabd01/socialblog-be/internal/middleware"
)

func PostRoutes(r *gin.Engine) {
	r.GET("/posts", postcontroller.Index)
	r.POST("/posts", middleware.AuthMiddleware(), postcontroller.Create)
	r.GET("/posts/:id", postcontroller.ShowDetail)
	r.PUT("/posts/:id", postcontroller.Update)
}
