package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
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
		"toType":     GetAttributeType,
		"toJson":     GetAttributeJson,
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

func Snakeze(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetAttributeName(in string) string {

	attribute := strings.Split(in, ":")
	attributeName := Camelize(attribute[0])

	return attributeName
}

func GetAttributeType(in string) string {
	attribute := strings.Split(in, ":")
	return attribute[1]
}

func GetAttributeJson(in string) string {
	name := GetAttributeName(in)
	name = Snakeze(name)

	return fmt.Sprintf("`json:\"%s\" binding:\"required\"`", name)
}
