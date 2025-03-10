package routes

import (
	"github.com/gin-gonic/gin"
	useractioncontroller "github.com/rianabd01/socialblog-be/internal/controllers/useraction-controller"
	"github.com/rianabd01/socialblog-be/internal/middleware"
)

func UserActionRoutes(r *gin.Engine) {
	r.POST("/api/like", middleware.AuthMiddleware(), useractioncontroller.LikeContent)
	r.POST("/api/unlike", middleware.AuthMiddleware(), useractioncontroller.UnlikeContent)
}
