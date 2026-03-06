package plugin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	registryURL = "https://raw.githubusercontent.com/dayvsonlima/catuaba-plugins/main/registry.yaml"
	cacheTTL    = 24 * time.Hour
)

// Registry holds the list of available plugins.
type Registry struct {
	Plugins map[string]RegistryEntry `yaml:"plugins"`
}

// RegistryEntry is a single plugin entry in the registry.
type RegistryEntry struct {
	Repository  string `yaml:"repository"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// LoadRegistry fetches the plugin registry, using a local cache with 24h TTL.
// If refresh is true, the cache is bypassed.
func LoadRegistry(refresh bool) (*Registry, error) {
	cachePath := registryCachePath()

	// Check cache unless refreshing
	if !refresh {
		if reg, err := loadCachedRegistry(cachePath); err == nil {
			return reg, nil
		}
	}

	// Fetch from remote
	reg, err := fetchRemoteRegistry()
	if err != nil {
		// Fall back to cache if available
		if cached, cacheErr := loadCachedRegistryIgnoreTTL(cachePath); cacheErr == nil {
			return cached, nil
		}
		return nil, fmt.Errorf("fetching registry: %w", err)
	}

	// Save to cache
	_ = saveRegistryCache(cachePath, reg)

	return reg, nil
}

func registryCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".catuaba", "registry.yaml")
}

func loadCachedRegistry(path string) (*Registry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Check TTL
	if time.Since(info.ModTime()) > cacheTTL {
		return nil, fmt.Errorf("cache expired")
	}

	return parseRegistryFile(path)
}

func loadCachedRegistryIgnoreTTL(path string) (*Registry, error) {
	return parseRegistryFile(path)
}

func parseRegistryFile(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}

	return &reg, nil
}

func fetchRemoteRegistry() (*Registry, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}

	return &reg, nil
}

func saveRegistryCache(path string, reg *Registry) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(reg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
