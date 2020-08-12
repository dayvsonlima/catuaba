package templates

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"
	"unicode"
)

func Render(fileName string, data interface{}) string {
	tmpl, _ := ioutil.ReadFile("cli/model/template.go.tmpl")

	t, _ := template.New("tmpm").Funcs(template.FuncMap{
		"toModelName": func(text string) string {
			return Camelize(text)
		},
		"toAttrName": GetAttributeName,
	}).Parse(string(tmpl))

	buf := &bytes.Buffer{}
	t.Execute(buf, data)

	return buf.String()
}

func Camelize(in string) string {
	runes := []rune(in)
	var out []rune

	for i, r := range runes {
		if r == '_' {
			continue
		}
		if i == 0 || runes[i-1] == '_' {
			out = append(out, unicode.ToUpper(r))
			continue
		}
		out = append(out, r)
	}

	return string(out)
}

func GetAttributeName(in string) string {

	attribute := strings.Split(in, ":")
	attributeName := Camelize(attribute[0])

	return attributeName
}
