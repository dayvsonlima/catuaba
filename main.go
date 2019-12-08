package main

import "github.com/gin-gonic/gin"


func ping(response *gin.Context) {
  response.JSON(200, gin.H{ "message": "pongozila" })
}

func main() {
  route := gin.Default()

  route.GET("/ping", ping)

  route.Run()
}
