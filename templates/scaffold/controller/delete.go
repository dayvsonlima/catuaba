package controller

var Delete = `package {{.Name|toPlural|toSnake}}

import (
	"application/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete .
func Delete(ctx *gin.Context) {
	var {{.Name|toVarName}} models.{{.Name|toModelName}}
  db.First(&{{.Name|toVarName}}, ctx.Param("id"))

	if {{.Name|toVarName}}.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&{{.Name|toVarName}})

	ctx.JSON(http.StatusOK, gin.H{"data": true})
}`
