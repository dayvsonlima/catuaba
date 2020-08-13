package main

import (
	"application/config"
	"application/database"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	config.DrawRoutes(engine)

	database.Migrations()
	defer database.Connection.Close()

	engine.Use(gin.Logger())
	engine.Run(":8080")
}
