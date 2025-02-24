package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rianabd01/socialblog-be/internal/controllers"
	"github.com/rianabd01/socialblog-be/internal/middleware"
)

func BlogRoutes(r *gin.Engine) {
	r.GET("/api/blogs", controller.Index)
	r.GET("/api/blogs/:id", controller.ShowDetail)
	r.POST("/api/blogs", middleware.AuthMiddleware(), controller.Create)
	r.PUT("/api/blogs/:id", middleware.AuthMiddleware(), controller.Update)
}
