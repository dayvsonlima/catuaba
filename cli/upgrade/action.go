package upgrade

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/dayvsonlima/catuaba/cli/output"
	"github.com/urfave/cli/v2"
)

const (
	githubReleaseURL = "https://api.github.com/repos/dayvsonlima/catuaba/releases/latest"
	installPkg       = "github.com/dayvsonlima/catuaba@latest"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func Action(c *cli.Context) error {
	currentVersion, _ := c.App.Metadata["version"].(string)
	if currentVersion == "" {
		return fmt.Errorf("could not determine current version")
	}

	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go is not installed or not in PATH")
	}

	output.Info("Current version: %s", currentVersion)
	output.Info("Checking for updates...")

	latest, err := fetchLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	cmp := compareVersions(latest, currentVersion)
	if cmp <= 0 {
		output.Success("Already up to date! (%s)", currentVersion)
		return nil
	}

	output.Info("New version available: %s → %s", currentVersion, latest)
	output.Info("Upgrading...")

	cmd := exec.Command("go", "install", installPkg)
	cmd.Stdout = c.App.Writer
	cmd.Stderr = c.App.ErrWriter
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	// Verify the upgrade
	out, err := exec.Command("catuaba", "--version").CombinedOutput()
	if err == nil {
		output.Success("Upgraded to %s", strings.TrimSpace(string(out)))
	} else {
		output.Success("Upgrade complete! New version: %s", latest)
	}

	return nil
}

func fetchLatestVersion() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(githubReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	version := strings.TrimPrefix(release.TagName, "v")
	if version == "" {
		return "", fmt.Errorf("empty version tag from GitHub")
	}

	return version, nil
}

// compareVersions compares two semver strings (e.g. "0.1.5" vs "0.1.4").
// Returns 1 if a > b, -1 if a < b, 0 if equal.
func compareVersions(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var va, vb int
		if i < len(partsA) {
			va, _ = strconv.Atoi(partsA[i])
		}
		if i < len(partsB) {
			vb, _ = strconv.Atoi(partsB[i])
		}
		if va > vb {
			return 1
		}
		if va < vb {
			return -1
		}
	}
	return 0
}
