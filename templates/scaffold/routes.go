package scaffold

var Routes = `

routes.GET("/{{.Name|toPlural|toSnake}}", {{.Name|toPlural|toSnake}}.Index)
routes.POST("/{{.Name|toPlural|toSnake}}", {{.Name|toPlural|toSnake}}.Create)
routes.GET("/{{.Name|toPlural|toSnake}}/:id", {{.Name|toPlural|toSnake}}.Show)
routes.PUT("/{{.Name|toPlural|toSnake}}/:id", {{.Name|toPlural|toSnake}}.Update)
routes.PATCH("/{{.Name|toPlural|toSnake}}/:id", {{.Name|toPlural|toSnake}}.Update)
routes.DELETE("/{{.Name|toPlural|toSnake}}/:id", {{.Name|toPlural|toSnake}}.Delete)
`
