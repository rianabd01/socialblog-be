package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rianabd01/socialblog-be/internal/controllers/usercontroller"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/signup", usercontroller.Signup)
	r.POST("/login", usercontroller.Login)
	r.GET("/auth/google/login", usercontroller.GoogleLogin)
	r.GET("/auth/google/callback", usercontroller.GoogleCallback)
}
