package controller

var Create = `package {{.Name|toPlural|toSnake}}

import (
	"application/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateParams struct {
	{{ range .Params}}{{. | toAttrName}} {{. | toType}} {{. | toJson}}
	{{ end }}
}

// Create .
func Create(ctx *gin.Context) {
	var params CreateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.Name|toVarName}} := &models.{{.Name|toModelName}}{
		{{ range .Params}}{{. | toAttrName}}: params.{{. | toAttrName}},
		{{ end }}
	}

	db.Create({{.Name|toVarName}})
	ctx.JSON(http.StatusOK, {{.Name|toVarName}})
}`
