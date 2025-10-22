package self

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/joeblew999/goup-util/pkg/self/output"
)

// Doctor validates the installation and all dependencies
func Doctor() error {
	result := output.DoctorResult{
		Installations: []output.InstallationInfo{},
		Issues:        []string{},
		Suggestions:   []string{},
	}

	// Check goup-util itself - look for ALL installations
	installations := findAllInstallations()

	if len(installations) == 0 {
		result.Issues = append(result.Issues, "goup-util not found in PATH")
		result.Suggestions = append(result.Suggestions, "Run: goup-util self setup")
	} else {
		for i, path := range installations {
			info := output.InstallationInfo{
				Path:     path,
				Active:   i == 0,
				Shadowed: i > 0,
			}
			result.Installations = append(result.Installations, info)
		}

		if len(installations) > 1 {
			result.Issues = append(result.Issues, "Multiple goup-util installations found")
			for i, path := range installations {
				if i > 0 {
					result.Suggestions = append(result.Suggestions, "Remove: "+path)
				}
			}
		}
	}

	// Check platform-specific package manager
	switch runtime.GOOS {
	case "darwin":
		if err := checkCommand("brew", "--version"); err != nil {
			result.Issues = append(result.Issues, "Homebrew not installed")
		}
	case "windows":
		if err := checkCommand("winget", "--version"); err != nil {
			result.Issues = append(result.Issues, "winget not found (optional)")
		}
	}

	// Check git
	if err := checkCommand("git", "--version"); err != nil {
		result.Issues = append(result.Issues, "git not installed")
		result.Suggestions = append(result.Suggestions, "Install git")
	}

	// Check go
	if err := checkCommand("go", "version"); err != nil {
		result.Issues = append(result.Issues, "go not installed")
		result.Suggestions = append(result.Suggestions, "Install go")
	}

	// Check task
	if err := checkCommand("task", "--version"); err != nil {
		result.Issues = append(result.Issues, "task not installed")
		result.Suggestions = append(result.Suggestions, "Install task")
	}

	output.OK("self doctor", result)
	return nil
}

// checkCommand checks if a command exists and runs successfully
func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// findAllInstallations finds all goup-util binaries in PATH
func findAllInstallations() []string {
	var installations []string

	// Get PATH
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	// Check each directory in PATH
	for _, dir := range paths {
		binaryPath := filepath.Join(dir, BinaryName)

		// Check if file exists and is executable
		if info, err := os.Stat(binaryPath); err == nil && !info.IsDir() {
			// Check if executable
			if info.Mode()&0111 != 0 {
				installations = append(installations, binaryPath)
			}
		}
	}

	return installations
}
