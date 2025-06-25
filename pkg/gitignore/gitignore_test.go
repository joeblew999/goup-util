package gitignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test with non-existent .gitignore
	gi := New(tempDir)
	if gi.Exists {
		t.Error("Expected .gitignore to not exist")
	}
	if gi.ProjectPath != tempDir {
		t.Errorf("Expected project path %s, got %s", tempDir, gi.ProjectPath)
	}

	// Create a .gitignore file
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("*.log\n# Comment\n\n.bin/\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test with existing .gitignore
	gi2 := New(tempDir)
	if !gi2.Exists {
		t.Error("Expected .gitignore to exist")
	}
}

func TestLoad(t *testing.T) {
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")

	content := "*.log\n# Comment\n\n.bin/\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	gi := New(tempDir)
	if err := gi.Load(); err != nil {
		t.Fatal(err)
	}

	expectedLines := []string{"*.log", "# Comment", "", ".bin/"}
	if len(gi.Lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(gi.Lines))
	}

	for i, expected := range expectedLines {
		if i < len(gi.Lines) && gi.Lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, gi.Lines[i])
		}
	}
}

func TestHasPattern(t *testing.T) {
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")

	content := "*.log\n.bin/\n# Comment\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	gi := New(tempDir)
	if err := gi.Load(); err != nil {
		t.Fatal(err)
	}

	if !gi.HasPattern("*.log") {
		t.Error("Expected to find pattern '*.log'")
	}

	if !gi.HasPattern(".bin/") {
		t.Error("Expected to find pattern '.bin/'")
	}

	if gi.HasPattern("nonexistent") {
		t.Error("Did not expect to find pattern 'nonexistent'")
	}
}

func TestInfo(t *testing.T) {
	tempDir := t.TempDir()

	// Test with non-existent .gitignore
	gi := New(tempDir)
	if err := gi.Load(); err != nil {
		t.Fatal(err)
	}

	info := gi.Info()
	if info["exists"].(bool) {
		t.Error("Expected exists to be false")
	}

	// Test with existing .gitignore
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	content := "*.log\n# Comment\n\n.bin/\n"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	gi2 := New(tempDir)
	if err := gi2.Load(); err != nil {
		t.Fatal(err)
	}

	info2 := gi2.Info()
	if !info2["exists"].(bool) {
		t.Error("Expected exists to be true")
	}
	if info2["lines"].(int) != 4 {
		t.Errorf("Expected 4 lines, got %d", info2["lines"].(int))
	}
	if info2["patterns"].(int) != 2 {
		t.Errorf("Expected 2 patterns, got %d", info2["patterns"].(int))
	}
	if info2["comments"].(int) != 1 {
		t.Errorf("Expected 1 comment, got %d", info2["comments"].(int))
	}
}
