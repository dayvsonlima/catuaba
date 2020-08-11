package model

const Template = `
package models

import "github.com/jinzhu/gorm"

// {{toModelName .Name}} model
type {{toModelName .Name}} struct {
	gorm.Model 
	Title     
	Content   
}
`
