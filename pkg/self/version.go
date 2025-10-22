package self

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joeblew999/goup-util/pkg/self/output"
)

// Version is set by the build process
var Version = "dev"

// ShowVersion displays the current version of goup-util
func ShowVersion() error {
	location := ""
	if path, err := exec.LookPath("goup-util"); err == nil {
		location = path
	}

	result := output.VersionResult{
		Version:  Version,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Location: location,
	}

	output.Print(result, "self version")
	return nil
}

// ShowStatus checks installation status and available updates
func ShowStatus() error {
	result := output.StatusResult{}

	// Check if installed
	installPath, err := exec.LookPath("goup-util")
	if err != nil {
		result.Installed = false
		result.UpdateAvailable = false
		output.Print(result, "self status")
		return nil
	}

	result.Installed = true
	result.CurrentVersion = normalizeVersion(Version)
	result.Location = installPath

	// Check for updates (from GitHub)
	latest, err := getLatestVersion(FullRepoName)
	if err == nil && latest != "" {
		result.LatestVersion = latest
		result.UpdateAvailable = (result.CurrentVersion != latest)
	} else {
		result.LatestVersion = ""
		result.UpdateAvailable = false
	}

	output.Print(result, "self status")
	return nil
}

// getLatestVersion fetches the latest release tag from GitHub
func getLatestVersion(repo string) (string, error) {
	cmd := exec.Command("git", "ls-remote", "--tags", "--refs",
		fmt.Sprintf("https://github.com/%s.git", repo))

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse output to find latest semver tag
	lines := strings.Split(string(output), "\n")
	var latest string

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		tag := filepath.Base(parts[1]) // Extract tag name from refs/tags/v1.2.3
		if strings.HasPrefix(tag, "v") && strings.Contains(tag, ".") {
			// Simple comparison - just use last tag (sorted by git)
			latest = tag
		}
	}

	return latest, nil
}
