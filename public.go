package Geopackage

import (
	"github.com/gin-gonic/gin"
	"github.com/khanjumed/geopackage/internal/routes"
)

func RegisterAllRoutes(r *gin.Engine) {
	routes.Register(r)
}
