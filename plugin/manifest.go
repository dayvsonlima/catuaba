package plugin

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Manifest represents the declarative plugin.yaml file.
type Manifest struct {
	Name        string     `yaml:"name"`
	Version     string     `yaml:"version"`
	Description string     `yaml:"description"`
	Author      string     `yaml:"author"`
	Repository  string     `yaml:"repository"`
	Variables   []Variable `yaml:"variables"`
	Dependencies []string  `yaml:"dependencies"`
	Directories []string   `yaml:"directories"`
	Files       []File     `yaml:"files"`
	Injections  []Injection `yaml:"injections"`
	EnvVars     []EnvVar   `yaml:"env_vars"`
	PostInstall []string   `yaml:"post_install"`
}

// Variable is a configurable plugin variable with a default value.
type Variable struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
}

// File maps a template to an output path.
type File struct {
	Template string `yaml:"template"`
	Output   string `yaml:"output"`
}

// Injection describes a code modification to an existing file.
type Injection struct {
	File      string `yaml:"file"`
	Action    string `yaml:"action"`
	Import    string `yaml:"import"`
	Method    string `yaml:"method"`
	Attribute string `yaml:"attribute"`
	Content   string `yaml:"content"`
	After     string `yaml:"after"`
}

// EnvVar is an environment variable to append to .env files.
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// LoadManifest reads and parses a plugin.yaml from the given path.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}

	return &m, nil
}
