package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Matches: templ ComponentName(params...) {
	reTemplFunc = regexp.MustCompile(`^templ\s+(\w+)\((.*?)\)\s*\{`)
	// Matches: type TypeName struct {
	reTypeStruct = regexp.MustCompile(`^type\s+(\w+)\s+struct\s*\{`)
)

// ParseComponents scans .templ files in componentsDir and extracts
// component signatures, whether they accept children, and Go types.
func ParseComponents(componentsDir string) (*ComponentsResult, error) {
	entries, err := os.ReadDir(componentsDir)
	if err != nil {
		return &ComponentsResult{}, nil
	}

	var components []ComponentInfo
	var types []ComponentType

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".templ") {
			continue
		}

		filePath := filepath.Join(componentsDir, e.Name())
		comps, typs, err := parseTemplFile(filePath, e.Name())
		if err != nil {
			continue
		}

		components = append(components, comps...)
		types = append(types, typs...)
	}

	return &ComponentsResult{
		Components: components,
		Types:      types,
	}, nil
}

func parseTemplFile(filePath, fileName string) ([]ComponentInfo, []ComponentType, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var components []ComponentInfo
	var types []ComponentType

	scanner := bufio.NewScanner(f)
	var currentComponent *ComponentInfo
	braceDepth := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check for templ function declaration
		if m := reTemplFunc.FindStringSubmatch(trimmed); m != nil {
			comp := ComponentInfo{
				Name:   m[1],
				File:   fileName,
				Params: parseParams(m[2]),
			}
			currentComponent = &comp
			braceDepth = 1
			continue
		}

		// Track braces to know when we're inside a component body
		if currentComponent != nil {
			braceDepth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")

			// Check for children slot
			if strings.Contains(trimmed, "{ children... }") {
				currentComponent.Children = true
			}

			if braceDepth <= 0 {
				components = append(components, *currentComponent)
				currentComponent = nil
				braceDepth = 0
			}
			continue
		}

		// Check for type struct declaration (outside of templ blocks)
		if m := reTypeStruct.FindStringSubmatch(trimmed); m != nil {
			typeName := m[1]
			fields := parseStructFields(scanner)
			types = append(types, ComponentType{
				Name:   typeName,
				Fields: fields,
			})
		}
	}

	// If we ended while still inside a component, add it
	if currentComponent != nil {
		components = append(components, *currentComponent)
	}

	return components, types, scanner.Err()
}

// parseParams parses a templ function parameter list like "label string, options []SelectOption"
func parseParams(raw string) []ComponentParam {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var params []ComponentParam
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// Split on last space to handle types like "[]SelectOption"
		idx := strings.LastIndex(part, " ")
		if idx < 0 {
			continue
		}
		name := strings.TrimSpace(part[:idx])
		typ := strings.TrimSpace(part[idx+1:])
		params = append(params, ComponentParam{Name: name, Type: typ})
	}
	return params
}

// parseStructFields reads struct fields until the closing brace
func parseStructFields(scanner *bufio.Scanner) []FieldInfo {
	var fields []FieldInfo
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "}" {
			break
		}
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fields = append(fields, FieldInfo{
				Name: parts[0],
				Type: parts[1],
			})
		}
	}
	return fields
}
