package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

)

// ParseModels parses all .go files in the models directory and extracts struct definitions.
func ParseModels(modelsDir string) ([]ModelInfo, error) {
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, err
	}

	var models []ModelInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(modelsDir, e.Name())
		parsed, err := parseModelFile(filePath)
		if err != nil {
			continue // skip unparseable files
		}
		models = append(models, parsed...)
	}

	return models, nil
}

func parseModelFile(filePath string) ([]ModelInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var models []ModelInfo
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			model := ModelInfo{
				Name:   typeSpec.Name.Name,
				Fields: parseFields(structType),
			}
			models = append(models, model)
		}
	}

	return models, nil
}

func parseFields(s *ast.StructType) []FieldInfo {
	var fields []FieldInfo
	for _, field := range s.Fields.List {
		fi := FieldInfo{}

		// Embedded type (e.g. gorm.Model)
		if len(field.Names) == 0 {
			fi.Name = typeString(field.Type)
			fi.Type = "embedded"
		} else {
			fi.Name = field.Names[0].Name
			fi.Type = typeString(field.Type)
		}

		// Parse struct tags
		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			fi.JSON = extractTag(tag, "json")
			fi.GORM = extractTag(tag, "gorm")
		}

		fields = append(fields, fi)
	}
	return fields
}

func typeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.ArrayType:
		return "[]" + typeString(t.Elt)
	case *ast.MapType:
		return "map[" + typeString(t.Key) + "]" + typeString(t.Value)
	default:
		return "unknown"
	}
}

func extractTag(tag, key string) string {
	// Find key:"value" in the tag string
	search := key + `:`
	idx := strings.Index(tag, search)
	if idx == -1 {
		return ""
	}

	rest := tag[idx+len(search):]
	if len(rest) == 0 || rest[0] != '"' {
		return ""
	}
	rest = rest[1:]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return ""
	}
	return rest[:end]
}
