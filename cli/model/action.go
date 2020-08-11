package model

import (
	"bytes"
	"fmt"
	"text/template"
	"unicode"

	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {

	t, _ := template.New("new_tmpl").Funcs(template.FuncMap{
		"toModelName": func(text string) string {
			return Camelize(text)
		},
	}).Parse(Template)

	buf := &bytes.Buffer{}

	data := struct {
		Name string
	}{
		Name: c.Args().Get(0),
	}

	t.Execute(buf, data)
	fmt.Println(buf)

	return nil
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
