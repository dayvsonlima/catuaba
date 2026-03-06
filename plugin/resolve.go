package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveSource determines the plugin source type and returns the local path
// to the plugin directory (downloading if needed).
//
// Supported sources:
//   - Local path: ./local-plugin or /absolute/path
//   - Git URL: github.com/user/repo
//   - Short name: "auth" (resolved via registry)
func ResolveSource(source string) (pluginDir string, cleanup func(), err error) {
	// Local path
	if strings.HasPrefix(source, ".") || strings.HasPrefix(source, "/") {
		absPath, err := filepath.Abs(source)
		if err != nil {
			return "", nil, fmt.Errorf("resolving path: %w", err)
		}

		if _, err := os.Stat(filepath.Join(absPath, "plugin.yaml")); os.IsNotExist(err) {
			return "", nil, fmt.Errorf("no plugin.yaml found at %s", absPath)
		}

		return absPath, func() {}, nil
	}

	// Git URL (contains "/" — e.g. github.com/user/repo)
	if strings.Contains(source, "/") {
		dir, cleanup, err := cloneRepo(source)
		if err != nil {
			return "", nil, err
		}
		return dir, cleanup, nil
	}

	// Short name — resolve via registry
	reg, err := LoadRegistry(false)
	if err != nil {
		return "", nil, fmt.Errorf("loading registry: %w", err)
	}

	entry, ok := reg.Plugins[source]
	if !ok {
		return "", nil, fmt.Errorf("plugin %q not found in registry. Use 'catuaba plugin list' to see available plugins", source)
	}

	dir, cleanup, err := cloneRepo(entry.Repository)
	if err != nil {
		return "", nil, err
	}
	return dir, cleanup, nil
}

// cloneRepo clones a git repository to a temp directory and returns the path.
func cloneRepo(repo string) (string, func(), error) {
	url := repo
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	tmpDir, err := os.MkdirTemp("", "catuaba-plugin-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	cmd := exec.Command("git", "clone", "--depth=1", url, tmpDir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning %s: %w", repo, err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "plugin.yaml")); os.IsNotExist(err) {
		cleanup()
		return "", nil, fmt.Errorf("cloned repo %s does not contain plugin.yaml", repo)
	}

	return tmpDir, cleanup, nil
}
