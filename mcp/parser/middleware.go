package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

)

// ParseMiddleware scans the middleware directory for Go files and extracts
// exported function names.
func ParseMiddleware(middlewareDir string) ([]MiddlewareInfo, error) {
	entries, err := os.ReadDir(middlewareDir)
	if err != nil {
		return nil, err
	}

	var middlewares []MiddlewareInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(middlewareDir, e.Name())
		funcs, err := parseExportedFuncs(filePath)
		if err != nil {
			continue
		}

		for _, name := range funcs {
			middlewares = append(middlewares, MiddlewareInfo{
				Name: name,
				File: e.Name(),
			})
		}
	}

	return middlewares, nil
}

func parseExportedFuncs(filePath string) ([]string, error) {
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
		if !ast.IsExported(funcDecl.Name.Name) {
			continue
		}
		names = append(names, funcDecl.Name.Name)
	}

	return names, nil
}
