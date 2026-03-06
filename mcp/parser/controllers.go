package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

)

// ParseControllers scans the controllers directory for Go files and extracts
// exported function names grouped by package (subdirectory).
func ParseControllers(controllersDir string) ([]ControllerInfo, error) {
	entries, err := os.ReadDir(controllersDir)
	if err != nil {
		return nil, err
	}

	var controllers []ControllerInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		pkgDir := filepath.Join(controllersDir, e.Name())
		funcs, err := parseControllerPackage(pkgDir)
		if err != nil || len(funcs) == 0 {
			continue
		}

		controllers = append(controllers, ControllerInfo{
			Package:   e.Name(),
			Functions: funcs,
		})
	}

	return controllers, nil
}

func parseControllerPackage(pkgDir string) ([]ControllerFunc, error) {
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, err
	}

	var funcs []ControllerFunc
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(pkgDir, e.Name())
		fileFuncs, err := parseExportedFunctions(filePath)
		if err != nil {
			continue
		}

		for _, name := range fileFuncs {
			funcs = append(funcs, ControllerFunc{
				Name: name,
				File: e.Name(),
			})
		}
	}

	return funcs, nil
}

func parseExportedFunctions(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// Only exported, non-method functions
		if funcDecl.Recv != nil {
			continue
		}
		if !ast.IsExported(funcDecl.Name.Name) {
			continue
		}
		names = append(names, funcDecl.Name.Name)
	}

	return names, nil
}
