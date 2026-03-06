package code_editor

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
	"strings"
)

// AddImport adds a new package import to Go source code using AST parsing.
// Falls back to regex-based approach if AST parsing fails.
func AddImport(code, path string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return addImportRegex(code, path)
	}

	// Check if import already exists
	quotedPath := strconv.Quote(path)
	for _, imp := range f.Imports {
		if imp.Path.Value == quotedPath {
			return code
		}
	}

	// Find the import declaration and add the new import
	added := false
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		newSpec := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: quotedPath,
			},
		}
		genDecl.Specs = append(genDecl.Specs, newSpec)
		added = true
		break
	}

	if !added {
		return addImportRegex(code, path)
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, f); err != nil {
		return addImportRegex(code, path)
	}

	return buf.String()
}

// AddImportIfNotExist adds an import only if it doesn't already exist in the code.
func AddImportIfNotExist(code, path string) string {
	if strings.Contains(code, path) {
		return code
	}
	return AddImport(code, path)
}

// AddAliasedImport adds a new aliased package import (e.g., alias "path/to/pkg").
func AddAliasedImport(code, alias, path string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return addAliasedImportRegex(code, alias, path)
	}

	quotedPath := strconv.Quote(path)
	for _, imp := range f.Imports {
		if imp.Path.Value == quotedPath {
			return code
		}
	}

	added := false
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		newSpec := &ast.ImportSpec{
			Name: ast.NewIdent(alias),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: quotedPath,
			},
		}
		genDecl.Specs = append(genDecl.Specs, newSpec)
		added = true
		break
	}

	if !added {
		return addAliasedImportRegex(code, alias, path)
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, f); err != nil {
		return addAliasedImportRegex(code, alias, path)
	}

	return buf.String()
}

func addAliasedImportRegex(code, alias, path string) string {
	idx := strings.LastIndex(code, ")")
	if idx == -1 {
		return code
	}
	return code[:idx] + "\t" + alias + " \"" + path + "\"\n" + code[idx:]
}

// addImportRegex is the legacy regex-based fallback for adding imports
func addImportRegex(code, path string) string {
	// Simple approach: find "import (" and add before closing ")"
	idx := strings.LastIndex(code, ")")
	if idx == -1 {
		return code
	}

	return code[:idx] + "\t\"" + path + "\"\n" + code[idx:]
}
