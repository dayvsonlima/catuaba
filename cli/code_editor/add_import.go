package code_editor

import (
	"regexp"
	"strings"
)

// AddImport
// Adds a new package to golang code inside the "import" statement
func AddImport(code, path string) string {

	rgx := `import(.|)\(((.|\n)+)\)`
	compiledRegex := regexp.MustCompile(rgx)
	submatch := compiledRegex.FindStringSubmatch(code)

	pkgString := submatch[2]
	pkgs := strings.Split(pkgString, "\n")

	pkgs = append(pkgs, path)
	normalized := normalizePkgs(pkgs)
	output := "import (\n" + strings.Join(normalized, "\n") + "\n)"

	return compiledRegex.ReplaceAllString(code, output)
}

func normalizePkgs(pkgs []string) []string {

	check := make(map[string]int)
	var output []string

	for _, value := range pkgs {
		current := strings.Trim(value, " ")
		check[current] = 1
	}

	for value, _ := range check {
		output = append(output, value)
	}

	return output
}
