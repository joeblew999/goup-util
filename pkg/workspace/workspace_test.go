package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFindWorkspace(t *testing.T) {
	tempDir := t.TempDir()

	// Create a go.work file in temp dir
	workContent := `go 1.21

use (
	./module1
	./module2
	./deep/module3
)`

	workPath := filepath.Join(tempDir, "go.work")
	if err := os.WriteFile(workPath, []byte(workContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test finding workspace by explicitly using filesystem traversal
	// (since go env GOWORK might find a real workspace)
	ws, err := findWorkspaceByTraversal(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ws.Exists {
		t.Error("Expected workspace to exist")
	}
	if len(ws.Modules) != 3 {
		t.Errorf("Expected 3 modules, got %d", len(ws.Modules))
	}

	expectedModules := []string{"./module1", "./module2", "./deep/module3"}
	for i, expected := range expectedModules {
		if i >= len(ws.Modules) || ws.Modules[i] != expected {
			t.Errorf("Expected module %d to be %s, got %s", i, expected, ws.Modules[i])
		}
	}
}

func TestFindWorkspaceInSubdirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create go.work in root
	workContent := `go 1.21

use ./module1`

	workPath := filepath.Join(tempDir, "go.work")
	if err := os.WriteFile(workPath, []byte(workContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create subdirectory
	subDir := filepath.Join(tempDir, "deep", "nested")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Search from subdirectory should find parent go.work
	// Use traversal method to avoid go env GOWORK interference
	ws, err := findWorkspaceByTraversal(subDir)
	if err != nil {
		t.Fatal(err)
	}
	if !ws.Exists {
		t.Error("Expected to find workspace in parent directory")
	}
	if ws.FilePath != workPath {
		t.Errorf("Expected workspace path %s, got %s", workPath, ws.FilePath)
	}
}

func TestHasModule(t *testing.T) {
	ws := &Workspace{
		Exists:  true,
		Modules: []string{"./module1", "./module2", "./deep/module3"},
	}

	if !ws.HasModule("./module1") {
		t.Error("Expected module1 to exist")
	}
	if !ws.HasModule("./deep/module3") {
		t.Error("Expected deep/module3 to exist")
	}
	if ws.HasModule("./nonexistent") {
		t.Error("Expected nonexistent module to not exist")
	}
}

func TestListModules(t *testing.T) {
	ws := &Workspace{
		Exists:  true,
		Modules: []string{"./module1", "./module2"},
	}

	modules := ws.ListModules()
	if len(modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(modules))
	}

	// Test that it returns a copy (mutations don't affect original)
	modules[0] = "modified"
	if ws.Modules[0] == "modified" {
		t.Error("ListModules should return a copy, not the original slice")
	}
}

func TestWorkspaceRoot(t *testing.T) {
	tempDir := t.TempDir()
	workPath := filepath.Join(tempDir, "go.work")

	ws := &Workspace{
		FilePath: workPath,
		Exists:   true,
	}

	root := ws.WorkspaceRoot()
	if root != tempDir {
		t.Errorf("Expected workspace root %s, got %s", tempDir, root)
	}

	// Test non-existent workspace
	ws2 := &Workspace{Exists: false}
	root2 := ws2.WorkspaceRoot()
	if root2 != "" {
		t.Errorf("Expected empty workspace root for non-existent workspace, got %s", root2)
	}
}

func TestInfo(t *testing.T) {
	// Test non-existent workspace
	ws := &Workspace{Exists: false}
	info := ws.Info()
	if info != "No go.work file found" {
		t.Errorf("Expected 'No go.work file found', got %s", info)
	}

	// Test existing workspace
	tempDir := t.TempDir()
	workPath := filepath.Join(tempDir, "go.work")

	ws2 := &Workspace{
		FilePath: workPath,
		Exists:   true,
		Modules:  []string{"./mod1", "./mod2"},
	}

	info2 := ws2.Info()
	expected := fmt.Sprintf("go.work: %s (2 modules)", workPath)
	if info2 != expected {
		t.Errorf("Expected %s, got %s", expected, info2)
	}
}

func TestLoadWorkspaceSingleUse(t *testing.T) {
	tempDir := t.TempDir()

	// Test single-line use statement
	workContent := `go 1.21

use ./single-module`

	workPath := filepath.Join(tempDir, "go.work")
	if err := os.WriteFile(workPath, []byte(workContent), 0644); err != nil {
		t.Fatal(err)
	}

	ws, err := loadWorkspace(workPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(ws.Modules) != 1 {
		t.Errorf("Expected 1 module, got %d", len(ws.Modules))
	}
	if ws.Modules[0] != "./single-module" {
		t.Errorf("Expected './single-module', got %s", ws.Modules[0])
	}
}
