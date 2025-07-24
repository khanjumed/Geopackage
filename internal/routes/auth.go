package routes

import (
	"github.com/khanjumed/geopackage/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}
}
