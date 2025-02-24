package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/controllers/blogcontroller"
	"github.com/rianabd01/socialblog-be/internal/middleware"
)

func BlogRoutes(r *gin.Engine) {
	r.GET("/api/blogs", blogcontroller.Index)
	r.GET("/api/blogs/:id", blogcontroller.ShowDetail)
	r.POST("/api/blogs", middleware.AuthMiddleware(), blogcontroller.Create)
	r.PUT("/api/blogs/:id", middleware.AuthMiddleware(), blogcontroller.Update)
}
