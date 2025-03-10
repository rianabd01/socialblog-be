package routes

import (
	"github.com/gin-gonic/gin"
	authcontroller "github.com/rianabd01/socialblog-be/internal/controllers/auth-controller"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/auth/signup", authcontroller.Signup)
	r.POST("/auth/login", authcontroller.Login)
	r.GET("/auth/google/login", authcontroller.GoogleLogin)
	r.GET("/auth/google/callback", authcontroller.GoogleCallback)
}
