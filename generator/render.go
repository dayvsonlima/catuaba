package generator

import (
	"bytes"
	"fmt"
	"text/template"
)

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"toModelName":    func(text string) string { return Camelize(text) },
		"toSnake":        Snakeze,
		"camelize":       Camelize,
		"toPlural":       Pluralize,
		"toLowerPlural":  LowerPlural,
		"toVarName":      CamelizeVar,
		"toAttrName":     GetAttributeName,
		"toType":         GetAttributeType,
		"toJson":         GetAttributeJson,
		"toJsonBinding":  GetAttributeJsonBinding,
		"toFormBinding":  GetFormBinding,
		"toFormLabel":    GetFormLabel,
		"toFormInput":    GetFormInputType,
		"toRawType":      GetAttributeRawType,
		"toSQLType":      GetSQLType,
		"toSQLDefault":   GetSQLDefault,
		"toSQLPrimaryKey": GetSQLPrimaryKey,
		"toSQLDefaultFor": GetSQLDefaultForDriver,
		"hasType":         HasType,
		"toColumnName":   GetSQLColumnName,
		"moduleName":     ModuleName,
	}
}

func Render(fileName string, data interface{}) (string, error) {
	content := LoadFile(fileName)

	t, err := template.New("tmpm").Funcs(templateFuncMap()).Parse(content)
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", fileName, err)
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", fileName, err)
	}

	return buf.String(), nil
}

func RenderFromContent(content string, data interface{}) (string, error) {
	t, err := template.New("tmpm").Funcs(templateFuncMap()).Parse(content)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}
