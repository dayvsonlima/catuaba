package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/dayvsonlima/catuaba/code_editor"
	"github.com/dayvsonlima/catuaba/generator"
)

// PluginData is the data available to plugin templates.
type PluginData struct {
	Name       string
	AppName    string
	ModuleName string
	Variables  map[string]string
}

// Install orchestrates the full plugin installation:
// dirs → files → injections → env vars → dependencies.
func Install(pluginDir string, m *Manifest, vars map[string]string) error {
	data := buildPluginData(m, vars)

	// 1. Create directories
	for _, dir := range m.Directories {
		if err := generator.Mkdir(dir); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	// 2. Render and write template files
	for _, f := range m.Files {
		tmplPath := filepath.Join(pluginDir, "templates", f.Template)
		content, err := os.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", f.Template, err)
		}

		rendered, err := renderPluginTemplate(string(content), data)
		if err != nil {
			return fmt.Errorf("rendering template %s: %w", f.Template, err)
		}

		outputPath := f.Output
		// Ensure parent directory exists
		dir := filepath.Dir(outputPath)
		if dir != "." {
			if err := generator.Mkdir(dir); err != nil {
				return fmt.Errorf("creating directory for %s: %w", outputPath, err)
			}
		}

		cwd, _ := os.Getwd()
		if err := os.WriteFile(filepath.Join(cwd, outputPath), []byte(rendered), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", outputPath, err)
		}
		output.Create(outputPath)
	}

	// 3. Apply injections
	for _, inj := range m.Injections {
		if err := applyInjection(inj, data); err != nil {
			return fmt.Errorf("injection on %s: %w", inj.File, err)
		}
	}

	// 4. Append env vars
	if len(m.EnvVars) > 0 {
		if err := appendEnvVars(m.EnvVars); err != nil {
			return fmt.Errorf("appending env vars: %w", err)
		}
	}

	// 5. Add Go dependencies
	if len(m.Dependencies) > 0 {
		if err := addDependencies(m.Dependencies); err != nil {
			return fmt.Errorf("adding dependencies: %w", err)
		}
	}

	return nil
}

func buildPluginData(m *Manifest, overrides map[string]string) PluginData {
	moduleName := generator.ModuleName()

	// Extract app name from module (last segment)
	appName := moduleName
	if idx := strings.LastIndex(moduleName, "/"); idx >= 0 {
		appName = moduleName[idx+1:]
	}

	// Build variables from manifest defaults, then apply overrides
	vars := make(map[string]string)
	for _, v := range m.Variables {
		vars[v.Name] = v.Default
	}
	for k, v := range overrides {
		vars[k] = v
	}

	return PluginData{
		Name:       m.Name,
		AppName:    appName,
		ModuleName: moduleName,
		Variables:  vars,
	}
}

func renderPluginTemplate(content string, data PluginData) (string, error) {
	funcMap := template.FuncMap{
		"toModelName":  generator.Camelize,
		"toSnake":      generator.Snakeze,
		"camelize":     generator.Camelize,
		"toPlural":     generator.Pluralize,
		"toLowerPlural": generator.LowerPlural,
		"toVarName":    generator.CamelizeVar,
		"toAttrName":   generator.GetAttributeName,
		"toType":       generator.GetAttributeType,
		"toJson":       generator.GetAttributeJson,
		"toJsonBinding": generator.GetAttributeJsonBinding,
		"moduleName":   generator.ModuleName,
	}

	t, err := template.New("plugin").Funcs(funcMap).Parse(content)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func applyInjection(inj Injection, data PluginData) error {
	// Render dynamic content in injection fields
	importPath := inj.Import
	content := inj.Content
	var err error

	if strings.Contains(importPath, "{{") {
		importPath, err = renderPluginTemplate(importPath, data)
		if err != nil {
			return fmt.Errorf("rendering import path: %w", err)
		}
	}
	if strings.Contains(content, "{{") {
		content, err = renderPluginTemplate(content, data)
		if err != nil {
			return fmt.Errorf("rendering injection content: %w", err)
		}
	}

	switch inj.Action {
	case "add_import":
		err = code_editor.EditFile(inj.File, func(code string) string {
			return code_editor.AddImportIfNotExist(code, importPath)
		})
		if err == nil {
			output.Inject(inj.File, "add import "+importPath)
		}

	case "insert_attribute":
		err = code_editor.EditFile(inj.File, func(code string) string {
			return code_editor.InsertAttribute(code, inj.Method, inj.Attribute)
		})
		if err == nil {
			output.Inject(inj.File, fmt.Sprintf("add %s to %s()", inj.Attribute, inj.Method))
		}

	case "insert_line_after":
		err = code_editor.EditFile(inj.File, func(code string) string {
			return InsertLineAfter(code, inj.After, content)
		})
		if err == nil {
			output.Inject(inj.File, "insert after "+inj.After)
		}

	case "insert_before_closing_brace":
		err = code_editor.EditFile(inj.File, func(code string) string {
			return InsertBeforeClosingBrace(code, content)
		})
		if err == nil {
			output.Inject(inj.File, "insert routes")
		}

	case "append":
		err = code_editor.EditFile(inj.File, func(code string) string {
			return AppendToFile(code, content)
		})
		if err == nil {
			output.Inject(inj.File, "append content")
		}

	default:
		return fmt.Errorf("unknown injection action %q", inj.Action)
	}

	return err
}

func appendEnvVars(envVars []EnvVar) error {
	for _, envFile := range []string{".env", ".env.example"} {
		cwd, _ := os.Getwd()
		path := filepath.Join(cwd, envFile)

		// Only append if the file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", envFile, err)
		}

		text := string(content)
		var newVars []string
		for _, ev := range envVars {
			line := ev.Name + "=" + ev.Value
			if !strings.Contains(text, ev.Name+"=") {
				newVars = append(newVars, line)
			}
		}

		if len(newVars) > 0 {
			text = AppendToFile(text, strings.Join(newVars, "\n"))
			if err := os.WriteFile(path, []byte(text), 0644); err != nil {
				return fmt.Errorf("writing %s: %w", envFile, err)
			}
			output.Inject(envFile, fmt.Sprintf("add %d env var(s)", len(newVars)))
		}
	}

	return nil
}

func addDependencies(deps []string) error {
	output.Info("Adding dependencies...")
	for _, dep := range deps {
		output.Info("  go get %s", dep)
	}
	return nil
}
