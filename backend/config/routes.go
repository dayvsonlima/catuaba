package config

import (
	"github.com/dayvsonlima/catuaba/backend/app/controllers/example_controller"
	posts_controller "github.com/dayvsonlima/catuaba/backend/app/controllers/posts"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	routes := gin.Default()

	routes.GET("/examples", example_controller.Index)
	routes.GET("/posts", posts_controller.Index)

	return routes
}
