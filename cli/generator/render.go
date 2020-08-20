package generator

import (
	"bytes"
	"text/template"
)

func Render(fileName string, data interface{}) string {
	content := LoadFile(fileName)

	t, _ := template.New("tmpm").Funcs(template.FuncMap{
		"toModelName": func(text string) string {
			return Camelize(text)
		},
		"toSnake":    Snakeze,
		"camelize":   Camelize,
		"toPlural":   Pluralize,
		"toVarName":  CamelizeVar,
		"toAttrName": GetAttributeName,
		"toType":     GetAttributeType,
		"toJson":     GetAttributeJson,
	}).Parse(content)

	buf := &bytes.Buffer{}
	t.Execute(buf, data)

	return buf.String()
}
