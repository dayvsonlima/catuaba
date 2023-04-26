package controller

var Show = `package {{.Name|toPlural|toSnake}}

import (
	"application/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show .
func Show(ctx *gin.Context) {
	var {{.Name|toVarName}} models.{{.Name|toModelName}}
	db.First(&{{.Name|toVarName}}, ctx.Param("id"))

	if {{.Name|toVarName}}.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "{{.Name|toVarName}} not found!",
		})

		return
	}

	ctx.JSON(http.StatusOK, {{.Name|toVarName}})
}`
