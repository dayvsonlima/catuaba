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
	if err := db.Where("id = ?", ctx.Param("id")).First(&{{.Name|toVarName}}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&{{.Name|toVarName}})

	ctx.JSON(http.StatusOK, gin.H{"data": true})
}`
