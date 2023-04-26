package templates

var Controller = `package {{.Name | toSnake }}

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// {{.MethodName | camelize}} .
func {{.MethodName | camelize }}(ctx *gin.Context) { }
`
