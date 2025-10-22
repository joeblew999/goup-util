package self

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	selfOutput "github.com/joeblew999/goup-util/pkg/self/output"
)

// ReleaseError outputs a JSON error for release command and exits
func ReleaseError(message string) error {
	selfOutput.PrintError("self release", fmt.Errorf("%s", message))
	return nil // Never reached because PrintError calls os.Exit
}

// Release performs the complete release process:
// - Runs tests
// - Builds all architectures
// - Creates git tag
// - Pushes to GitHub
// Release performs the complete release process:
// - Runs tests
// - Builds all architectures
// - Creates git tag
// - Pushes to GitHub
func Release(version string) error {
	result := selfOutput.ReleaseResult{
		TestsPassed: false,
		Built:       false,
		Tagged:      false,
		Pushed:      false,
		Binaries:    []string{},
	}

	// Validate and normalize version
	version = normalizeVersion(version)
	if err := validateVersion(version); err != nil {
		return err
	}
	result.Version = version

	// Check if working directory is clean
	if err := exec.Command("git", "diff-index", "--quiet", "HEAD", "--").Run(); err != nil {
		return fmt.Errorf("working directory is not clean. Please commit changes first")
	}

	// Run tests
	testCmd := exec.Command("go", "test", "./...")
	if _, err := testCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	// Run race tests
	raceCmd := exec.Command("go", "test", "-race", "./...")
	if _, err := raceCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("race tests failed: %w", err)
	}

	result.TestsPassed = true

	// Build self with obfuscation (use remote mode for releases)
	// Obfuscate releases for security
	if err := Build(BuildOptions{UseLocal: false, Obfuscate: true}); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Collect built binaries
	for _, arch := range SupportedArchitectures() {
		result.Binaries = append(result.Binaries, fmt.Sprintf("goup-util-%s", arch.Suffix))
	}
	result.Built = true

	// Create tag
	if err := exec.Command("git", "tag", "-a", version, "-m", "Release "+version).Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	result.Tagged = true

	// Push tag
	if err := exec.Command("git", "push", "origin", version).Run(); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}
	result.Pushed = true

	selfOutput.Print(result, "self release")
	return nil
}

// normalizeVersion handles version bumping (patch/minor/major) and v prefix
func normalizeVersion(version string) string {
	// Handle bump types
	if version == "patch" || version == "minor" || version == "major" {
		currentTag, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
		if err != nil {
			return "v1.0.0"
		}
		current := strings.TrimSpace(string(currentTag))
		return bumpVersion(current, version)
	}
	
	// Add v prefix if missing
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	
	return version
}

// validateVersion checks if version follows semantic versioning
func validateVersion(version string) error {
	if !regexp.MustCompile(`^v\d+\.\d+\.\d+$`).MatchString(version) {
		return fmt.Errorf("invalid version format: %s (use v1.2.3, patch, minor, or major)", version)
	}
	return nil
}

// bumpVersion increments version number based on bump type
func bumpVersion(current, bumpType string) string {
	current = strings.TrimPrefix(current, "v")
	parts := strings.Split(current, ".")
	if len(parts) != 3 {
		return "v1.0.0"
	}
	
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	
	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}
	
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}
