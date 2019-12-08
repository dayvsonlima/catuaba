package config

import "github.com/gin-gonic/gin"

import(
  "github.com/dayvsonlima/catuaba/app/controllers/example_controller"
)

func Routes() {
  routes := gin.Default()

  routes.GET("/", example_controller.Index)

  routes.Run()
}
