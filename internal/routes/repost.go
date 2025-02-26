package routes

import (
	"github.com/gin-gonic/gin"
	repostcontroller "github.com/rianabd01/socialblog-be/internal/controllers/repost-controller"
	"github.com/rianabd01/socialblog-be/internal/middleware"
)

func RepostRoutes(r *gin.Engine) {
	r.GET("/api/reposts", repostcontroller.Index)
	r.GET("/api/reposts/:id", repostcontroller.ShowDetail)
	r.POST("/api/reposts", middleware.AuthMiddleware(), repostcontroller.Create)
	r.PUT("/api/reposts/:id", middleware.AuthMiddleware(), repostcontroller.Update)
}
