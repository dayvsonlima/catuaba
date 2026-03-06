package plugin

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Tracker stores which plugins are installed in the project.
type Tracker struct {
	Installed map[string]TrackerEntry `yaml:"installed"`
}

// TrackerEntry records metadata about an installed plugin.
type TrackerEntry struct {
	Version      string   `yaml:"version"`
	InstalledAt  string   `yaml:"installed_at"`
	FilesCreated []string `yaml:"files_created"`
}

const trackerDir = ".catuaba"
const trackerFile = ".catuaba/plugins.yaml"

// LoadTracker reads the tracker file from the current project.
func LoadTracker() (*Tracker, error) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, trackerFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var t Tracker
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

// TrackInstall records a plugin installation in .catuaba/plugins.yaml.
func TrackInstall(m *Manifest) error {
	cwd, _ := os.Getwd()

	// Load existing tracker or create new
	tracker, err := LoadTracker()
	if err != nil {
		tracker = &Tracker{Installed: make(map[string]TrackerEntry)}
	}

	// Collect created files
	var files []string
	for _, f := range m.Files {
		files = append(files, f.Output)
	}

	tracker.Installed[m.Name] = TrackerEntry{
		Version:      m.Version,
		InstalledAt:  time.Now().Format(time.RFC3339),
		FilesCreated: files,
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Join(cwd, trackerDir), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(tracker)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(cwd, trackerFile), data, 0644)
}
