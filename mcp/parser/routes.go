package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

)

// ParseRoutes parses a config/routes.go file and extracts route definitions.
// It looks for method calls like r.GET("/path", handler) or api.POST("/path", handler).
func ParseRoutes(filePath string) ([]RouteInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var routes []RouteInfo
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		method := strings.ToUpper(sel.Sel.Name)
		if !isHTTPMethod(method) {
			return true
		}

		if len(call.Args) < 2 {
			return true
		}

		// First arg is the path
		pathLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || pathLit.Kind != token.STRING {
			return true
		}
		path := strings.Trim(pathLit.Value, `"`)

		// In Gin, the last arg is the handler; middle args are middleware
		handlerArgs := call.Args[1:]
		handler := exprToString(handlerArgs[len(handlerArgs)-1])

		var middleware []string
		for _, arg := range handlerArgs[:len(handlerArgs)-1] {
			mw := exprToString(arg)
			if mw != "" {
				middleware = append(middleware, mw)
			}
		}

		route := RouteInfo{
			Method:     method,
			Path:       path,
			Handler:    handler,
			Middleware: middleware,
		}
		routes = append(routes, route)

		return true
	})

	return routes, nil
}

func isHTTPMethod(m string) bool {
	switch m {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	}
	return false
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.Ident:
		return e.Name
	case *ast.CallExpr:
		return exprToString(e.Fun) + "()"
	default:
		return ""
	}
}
