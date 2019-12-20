package example_controller

import "github.com/gin-gonic/gin"

func Index(response *gin.Context) {
  response.JSON(200, gin.H{ "message": "this is one index page" })
}
