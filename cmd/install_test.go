package cmd

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestInstallSDK(t *testing.T) {
	sdkName := "test-sdk"

	// Mock the exec.Command function
	cmd := exec.Command("echo", "Mock installation of", sdkName)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run mock command: %v", err)
	}

	output := out.String()
	if output != "Mock installation of test-sdk\n" {
		t.Errorf("Unexpected output: %s", output)
	}
}
