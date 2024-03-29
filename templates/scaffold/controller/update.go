package controller

var Update = `package {{.Name|toPlural|toSnake}}

import (
	"net/http"

	"application/app/models"

	"github.com/gin-gonic/gin"
)

type UpdateParams struct {
	{{ range .Params}}{{. | toAttrName}} {{. | toType}} {{. | toJson}}
	{{ end }}
}

// Update .
func Update(ctx *gin.Context) {
	var {{.Name|toVarName}} models.{{.Name|toModelName}}
  db.First(&{{.Name|toVarName}}, ctx.Param("id"))

	if {{.Name|toVarName}}.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var params UpdateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&{{.Name|toVarName}}).Updates(params)

	ctx.JSON(http.StatusOK, {{.Name|toVarName}})
}`
