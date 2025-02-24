package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rianabd01/socialblog-be/internal/controllers"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/signup", controller.Signup)
	r.POST("/auth/login", controller.Login)
	r.GET("/auth/google/login", controller.GoogleLogin)
	r.GET("/auth/google/callback", controller.GoogleCallback)
}
