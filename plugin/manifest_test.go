package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
name: test-plugin
version: "1.0.0"
description: "A test plugin"
author: "tester"

variables:
  - name: foo
    default: "bar"

dependencies:
  - "github.com/some/dep@v1.0.0"

directories:
  - "app/test"

files:
  - template: "test.go.tmpl"
    output: "app/test/test.go"

injections:
  - file: "config/routes.go"
    action: "add_import"
    import: "test/pkg"

env_vars:
  - name: "TEST_VAR"
    value: "test_value"

post_install:
  - "Run go mod tidy"
`
	path := filepath.Join(dir, "plugin.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yamlContent), 0644))

	m, err := LoadManifest(path)
	require.NoError(t, err)

	assert.Equal(t, "test-plugin", m.Name)
	assert.Equal(t, "1.0.0", m.Version)
	assert.Equal(t, "A test plugin", m.Description)
	assert.Equal(t, "tester", m.Author)

	assert.Len(t, m.Variables, 1)
	assert.Equal(t, "foo", m.Variables[0].Name)
	assert.Equal(t, "bar", m.Variables[0].Default)

	assert.Len(t, m.Dependencies, 1)
	assert.Len(t, m.Directories, 1)
	assert.Len(t, m.Files, 1)
	assert.Equal(t, "test.go.tmpl", m.Files[0].Template)
	assert.Equal(t, "app/test/test.go", m.Files[0].Output)

	assert.Len(t, m.Injections, 1)
	assert.Equal(t, "add_import", m.Injections[0].Action)

	assert.Len(t, m.EnvVars, 1)
	assert.Len(t, m.PostInstall, 1)
}

func TestLoadManifest_NotFound(t *testing.T) {
	_, err := LoadManifest("/nonexistent/plugin.yaml")
	assert.Error(t, err)
}
