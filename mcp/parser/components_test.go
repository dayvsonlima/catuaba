package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseComponents(t *testing.T) {
	dir := t.TempDir()

	button := `package components

templ Button(label string, variant string) {
	<button class={ "btn btn-" + variant }>{ label }</button>
}
`
	card := `package components

templ Card(title string) {
	<div class="card">
		<h2>{ title }</h2>
		<div>
			{ children... }
		</div>
	</div>
}
`
	selectComp := `package components

type SelectOption struct {
	Value string
	Label string
}

templ Select(name string, options []SelectOption) {
	<select name={ name }>
		for _, opt := range options {
			<option value={ opt.Value }>{ opt.Label }</option>
		}
	</select>
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "button.templ"), []byte(button), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "card.templ"), []byte(card), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "select.templ"), []byte(selectComp), 0644))

	result, err := ParseComponents(dir)
	require.NoError(t, err)
	assert.Len(t, result.Components, 3)

	comps := make(map[string]ComponentInfo)
	for _, c := range result.Components {
		comps[c.Name] = c
	}

	// Button
	btn := comps["Button"]
	assert.Equal(t, "button.templ", btn.File)
	assert.Len(t, btn.Params, 2)
	assert.Equal(t, "label", btn.Params[0].Name)
	assert.Equal(t, "string", btn.Params[0].Type)
	assert.Equal(t, "variant", btn.Params[1].Name)
	assert.Equal(t, "string", btn.Params[1].Type)
	assert.False(t, btn.Children)

	// Card (with children)
	crd := comps["Card"]
	assert.Equal(t, "card.templ", crd.File)
	assert.Len(t, crd.Params, 1)
	assert.Equal(t, "title", crd.Params[0].Name)
	assert.True(t, crd.Children)

	// Select (with custom type)
	sel := comps["Select"]
	assert.Equal(t, "select.templ", sel.File)
	assert.Len(t, sel.Params, 2)
	assert.Equal(t, "options", sel.Params[1].Name)
	assert.Equal(t, "[]SelectOption", sel.Params[1].Type)
	assert.False(t, sel.Children)

	// Types
	assert.Len(t, result.Types, 1)
	assert.Equal(t, "SelectOption", result.Types[0].Name)
	assert.Len(t, result.Types[0].Fields, 2)
	assert.Equal(t, "Value", result.Types[0].Fields[0].Name)
	assert.Equal(t, "string", result.Types[0].Fields[0].Type)
}

func TestParseComponents_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := ParseComponents(dir)
	require.NoError(t, err)
	assert.Empty(t, result.Components)
	assert.Empty(t, result.Types)
}

func TestParseComponents_NoDir(t *testing.T) {
	result, err := ParseComponents("/nonexistent/path")
	require.NoError(t, err)
	assert.Empty(t, result.Components)
}

func TestParseComponents_NoParams(t *testing.T) {
	dir := t.TempDir()

	spinner := `package components

templ Spinner() {
	<div class="spinner"></div>
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "spinner.templ"), []byte(spinner), 0644))

	result, err := ParseComponents(dir)
	require.NoError(t, err)
	assert.Len(t, result.Components, 1)
	assert.Equal(t, "Spinner", result.Components[0].Name)
	assert.Empty(t, result.Components[0].Params)
	assert.False(t, result.Components[0].Children)
}
