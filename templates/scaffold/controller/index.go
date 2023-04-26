package controller

var Index = `package {{.Name|toPlural|toSnake}}

import (
	"application/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index .
func Index(ctx *gin.Context) {
	var {{.Name|toVarName|toPlural}} []models.{{.Name|toModelName}}
	db.Find(&{{.Name|toVarName|toPlural}})
	ctx.JSON(http.StatusOK, {{.Name|toVarName|toPlural}})
}`
