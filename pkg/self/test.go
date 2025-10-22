package self

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joeblew999/goup-util/pkg/self/output"
)

// TestBootstrap generates and tests bootstrap scripts locally
func TestBootstrap() error {
	result := output.TestResult{
		Phase:  "bootstrap_test",
		Passed: false,
		Steps:  []string{},
		Errors: []string{},
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to get current directory: %v", err))
		output.Print(result, "self test")
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Step 1: Build with local mode
	result.Steps = append(result.Steps, "Building with --local flag")
	if err := Build(BuildOptions{UseLocal: true}); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("build failed: %v", err))
		output.Print(result, "self test")
		return fmt.Errorf("build failed: %w", err)
	}

	// Step 2: Verify scripts exist
	result.Steps = append(result.Steps, "Verifying bootstrap scripts exist")
	scriptsDir := filepath.Join(currentDir, "scripts")
	macosScript := filepath.Join(scriptsDir, "macos-bootstrap.sh")
	windowsScript := filepath.Join(scriptsDir, "windows-bootstrap.ps1")

	if _, err := os.Stat(macosScript); err != nil {
		result.Errors = append(result.Errors, "macOS script not found")
		output.Print(result, "self test")
		return fmt.Errorf("macOS script not found: %w", err)
	}

	if _, err := os.Stat(windowsScript); err != nil {
		result.Errors = append(result.Errors, "Windows script not found")
		output.Print(result, "self test")
		return fmt.Errorf("Windows script not found: %w", err)
	}

	// Step 3: Validate script content
	result.Steps = append(result.Steps, "Validating script content")

	if err := validateScriptContains(macosScript, "LOCAL MODE"); err != nil {
		result.Errors = append(result.Errors, "macOS script missing LOCAL MODE marker")
		output.Print(result, "self test")
		return err
	}

	if err := validateScriptContains(macosScript, "self setup"); err != nil {
		result.Errors = append(result.Errors, "macOS script missing 'self setup' command")
		output.Print(result, "self test")
		return err
	}

	if err := validateScriptContains(windowsScript, "LOCAL MODE"); err != nil {
		result.Errors = append(result.Errors, "Windows script missing LOCAL MODE marker")
		output.Print(result, "self test")
		return err
	}

	if err := validateScriptContains(windowsScript, "self setup"); err != nil {
		result.Errors = append(result.Errors, "Windows script missing 'self setup' command")
		output.Print(result, "self test")
		return err
	}

	// Step 4: Test binary execution
	result.Steps = append(result.Steps, "Testing binary execution")
	arch := CurrentArchitecture()
	if arch == nil {
		result.Errors = append(result.Errors, "unsupported architecture")
		output.Print(result, "self test")
		return fmt.Errorf("unsupported architecture: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	binaryPath := filepath.Join(currentDir, arch.BinaryName())
	if _, err := os.Stat(binaryPath); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("binary not found: %s", binaryPath))
		output.Print(result, "self test")
		return fmt.Errorf("binary not found: %s", binaryPath)
	}

	cmd := exec.Command(binaryPath, "--help")
	if err := cmd.Run(); err != nil {
		result.Errors = append(result.Errors, "binary execution failed")
		output.Print(result, "self test")
		return fmt.Errorf("binary execution failed: %w", err)
	}

	// All tests passed!
	result.Passed = true
	output.Print(result, "self test")
	return nil
}

// validateScriptContains checks if a script file contains expected text
func validateScriptContains(scriptPath, expected string) error {
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", scriptPath, err)
	}

	if !strings.Contains(string(content), expected) {
		return fmt.Errorf("%s does not contain %q", filepath.Base(scriptPath), expected)
	}

	return nil
}
