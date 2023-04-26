package controller

var Shared = `package {{.Name|toPlural|toSnake}}

import "application/database"

var (
	db = database.ORM()
)`
