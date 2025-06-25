package workspace

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Workspace represents a Go workspace
type Workspace struct {
	FilePath string
	Exists   bool
	Modules  []string
}

// FindWorkspace searches for go.work file in current and parent directories
func FindWorkspace(startPath string) (*Workspace, error) {
	// First try go env GOWORK (most reliable)
	cmd := exec.Command("go", "env", "GOWORK")
	output, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		workFile := strings.TrimSpace(string(output))
		return loadWorkspace(workFile)
	}

	// Fallback to filesystem traversal
	currentDir := startPath
	if currentDir == "" {
		currentDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	for {
		workFile := filepath.Join(currentDir, "go.work")
		if _, err := os.Stat(workFile); err == nil {
			return loadWorkspace(workFile)
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // Reached root
			break
		}
		currentDir = parentDir
	}

	return &Workspace{Exists: false}, nil
}

// loadWorkspace reads and parses a go.work file
func loadWorkspace(filePath string) (*Workspace, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open go.work: %w", err)
	}
	defer file.Close()

	ws := &Workspace{
		FilePath: filePath,
		Exists:   true,
		Modules:  []string{},
	}

	scanner := bufio.NewScanner(file)
	inUseBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "use") {
			inUseBlock = true
			// Handle single-line use statements
			if strings.Contains(line, "(") {
				continue // Multi-line use block
			} else {
				// Single use statement: use ./path
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					ws.Modules = append(ws.Modules, parts[1])
				}
				continue
			}
		}

		if inUseBlock {
			if strings.Contains(line, ")") {
				inUseBlock = false
				continue
			}
			if line != "" && !strings.HasPrefix(line, "//") {
				ws.Modules = append(ws.Modules, line)
			}
		}
	}

	return ws, scanner.Err()
}

// HasModule checks if a module path is already in the workspace
func (w *Workspace) HasModule(modulePath string) bool {
	for _, mod := range w.Modules {
		if mod == modulePath {
			return true
		}
	}
	return false
}

// AddModule adds a module to the workspace (if not already present)
func (w *Workspace) AddModule(modulePath string, force bool) error {
	if !w.Exists {
		return fmt.Errorf("no go.work file found")
	}

	if w.HasModule(modulePath) {
		return nil // Already exists - idempotent
	}

	if !force {
		return fmt.Errorf("module %s not in workspace (use --force to add)", modulePath)
	}

	// Use go work use command (safer than manual file editing)
	cmd := exec.Command("go", "work", "use", modulePath)
	cmd.Dir = filepath.Dir(w.FilePath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add module to workspace: %w", err)
	}

	// Reload to get updated state
	updated, err := loadWorkspace(w.FilePath)
	if err != nil {
		return err
	}

	w.Modules = updated.Modules
	return nil
}

// RemoveModule removes a module from the workspace
func (w *Workspace) RemoveModule(modulePath string, force bool) error {
	if !w.Exists {
		return fmt.Errorf("no go.work file found")
	}

	if !w.HasModule(modulePath) {
		return nil // Already removed - idempotent
	}

	if !force {
		return fmt.Errorf("module %s exists in workspace (use --force to remove)", modulePath)
	}

	// Use go work drop command (safer than manual file editing)
	cmd := exec.Command("go", "work", "drop", modulePath)
	cmd.Dir = filepath.Dir(w.FilePath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove module from workspace: %w", err)
	}

	// Reload to get updated state
	updated, err := loadWorkspace(w.FilePath)
	if err != nil {
		return err
	}

	w.Modules = updated.Modules
	return nil
}

// Info returns information about the workspace
func (w *Workspace) Info() string {
	if !w.Exists {
		return "No go.work file found"
	}

	return fmt.Sprintf("go.work: %s (%d modules)", w.FilePath, len(w.Modules))
}

// String returns a summary
func (w *Workspace) String() string {
	if !w.Exists {
		return "No workspace"
	}
	return fmt.Sprintf("Workspace: %s", filepath.Base(filepath.Dir(w.FilePath)))
}

// ListModules returns all modules in the workspace
func (w *Workspace) ListModules() []string {
	if !w.Exists {
		return []string{}
	}
	return append([]string{}, w.Modules...) // Return copy
}

// WorkspaceRoot returns the directory containing the go.work file
func (w *Workspace) WorkspaceRoot() string {
	if !w.Exists {
		return ""
	}
	return filepath.Dir(w.FilePath)
}

// findWorkspaceByTraversal searches for go.work using only filesystem traversal
// This is useful for testing when we don't want go env GOWORK to interfere
func findWorkspaceByTraversal(startPath string) (*Workspace, error) {
	currentDir := startPath
	if currentDir == "" {
		var err error
		currentDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	for {
		workFile := filepath.Join(currentDir, "go.work")
		if _, err := os.Stat(workFile); err == nil {
			return loadWorkspace(workFile)
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // Reached root
			break
		}
		currentDir = parentDir
	}

	return &Workspace{Exists: false}, nil
}
