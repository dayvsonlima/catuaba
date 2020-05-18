package config

import (
	"github.com/dayvsonlima/catuaba/app/controllers/posts_controller"
	"github.com/gin-gonic/gin"
)

// Routes : put your routes here
func Routes() *gin.Engine {
	routes := gin.Default()

	routes.GET("/posts", posts_controller.Index)

	return routes
}
