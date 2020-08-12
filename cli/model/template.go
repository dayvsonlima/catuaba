package model

const Template = `package models

import "github.com/jinzhu/gorm"

// {{.Name | toModelName}} model
type {{.Name | toModelName}} struct {
	gorm.Model 
	{{ range .Params}}{{. | toAttrName}} {{. | toType}} {{. | toJson}}
	{{ end }}
}
`
