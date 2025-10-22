package self

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// InstallDeps installs system dependencies based on platform.
// For macOS: Homebrew, git, go, task
// For Windows: git, go, task via winget
// For Linux: git, go, task via package manager
func InstallDeps() error {
	switch runtime.GOOS {
	case "darwin":
		return installMacOSDeps()
	case "windows":
		return installWindowsDeps()
	case "linux":
		return installLinuxDeps()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// installMacOSDeps installs dependencies on macOS using Homebrew
func installMacOSDeps() error {
	// 1. Check/Install Homebrew
	if !commandExists("brew") {
		fmt.Println("📥 Homebrew not found. Installing...")
		cmd := exec.Command("/bin/bash", "-c",
			`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Homebrew: %w", err)
		}
		fmt.Println("✅ Homebrew installed")
	} else {
		fmt.Println("✅ Homebrew already installed")
	}

	// 2. Install required packages
	packages := []string{"git", "go", "go-task"}

	for _, pkg := range packages {
		if err := brewInstall(pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	return nil
}

// brewInstall installs a package via Homebrew if not already installed
func brewInstall(pkg string) error {
	// Check if already installed
	checkCmd := exec.Command("brew", "list", pkg)
	if err := checkCmd.Run(); err == nil {
		fmt.Printf("✅ %s already installed\n", pkg)
		return nil
	}

	// Install package
	fmt.Printf("📥 Installing %s via Homebrew...\n", pkg)
	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("brew install %s failed: %w", pkg, err)
	}

	fmt.Printf("✅ %s installed\n", pkg)
	return nil
}

// installWindowsDeps installs dependencies on Windows using winget
func installWindowsDeps() error {
	// Check for winget
	if !commandExists("winget") {
		return fmt.Errorf("winget not found. Please install App Installer from Microsoft Store: https://aka.ms/getwinget")
	}
	fmt.Println("✅ winget found")

	// Install required packages
	packages := []struct {
		id   string
		name string
	}{
		{"Git.Git", "Git"},
		{"GoLang.Go", "Go"},
		{"Task.Task", "Task"},
	}

	for _, pkg := range packages {
		if err := wingetInstall(pkg.id, pkg.name); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg.name, err)
		}
	}

	return nil
}

// wingetInstall installs a package via winget if not already installed
func wingetInstall(id, name string) error {
	// Check if already installed
	checkCmd := exec.Command("winget", "list", "--id", id, "--exact")
	var out bytes.Buffer
	checkCmd.Stdout = &out
	checkCmd.Stderr = &out
	if err := checkCmd.Run(); err == nil && strings.Contains(out.String(), id) {
		fmt.Printf("✅ %s already installed\n", name)
		return nil
	}

	// Install package
	fmt.Printf("📥 Installing %s via winget...\n", name)
	cmd := exec.Command("winget", "install", "--id", id, "--exact", "--silent", "--accept-source-agreements", "--accept-package-agreements")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("winget install %s failed: %w", id, err)
	}

	fmt.Printf("✅ %s installed\n", name)
	return nil
}

// installLinuxDeps installs dependencies on Linux
func installLinuxDeps() error {
	// Detect package manager
	var pkgManager string
	var installCmd []string

	if commandExists("apt-get") {
		pkgManager = "apt-get"
		installCmd = []string{"sudo", "apt-get", "install", "-y"}
	} else if commandExists("yum") {
		pkgManager = "yum"
		installCmd = []string{"sudo", "yum", "install", "-y"}
	} else if commandExists("dnf") {
		pkgManager = "dnf"
		installCmd = []string{"sudo", "dnf", "install", "-y"}
	} else if commandExists("pacman") {
		pkgManager = "pacman"
		installCmd = []string{"sudo", "pacman", "-S", "--noconfirm"}
	} else {
		return fmt.Errorf("no supported package manager found (apt-get, yum, dnf, pacman)")
	}

	fmt.Printf("✅ Using package manager: %s\n", pkgManager)

	// Install git and go
	packages := []string{"git", "golang"}
	for _, pkg := range packages {
		if commandExists(pkg) {
			fmt.Printf("✅ %s already installed\n", pkg)
			continue
		}

		fmt.Printf("📥 Installing %s via %s...\n", pkg, pkgManager)
		cmd := exec.Command(installCmd[0], append(installCmd[1:], pkg)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
		fmt.Printf("✅ %s installed\n", pkg)
	}

	// Install task via go install if not available in package manager
	if !commandExists("task") {
		fmt.Println("📥 Installing task via go install...")
		cmd := exec.Command("go", "install", "github.com/go-task/task/v3/cmd/task@latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install task: %w", err)
		}
		fmt.Println("✅ task installed")
	} else {
		fmt.Println("✅ task already installed")
	}

	return nil
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
