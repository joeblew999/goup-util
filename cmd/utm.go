package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var utmCmd = &cobra.Command{
	Use:   "utm",
	Short: "Control UTM virtual machines",
	Long: `Control UTM virtual machines using utmctl.

This command is a wrapper around utmctl for convenient VM automation.
Requires UTM to be installed and QEMU guest agent running in the VM.

Examples:
  # List all VMs
  goup-util utm list

  # Check VM status
  goup-util utm status "Windows 11"

  # Execute command in VM
  goup-util utm exec "Windows 11" -- build windows examples/hybrid-dashboard

  # Execute Task in VM
  goup-util utm task "Windows 11" build:hybrid:windows

  # Pull file from VM
  goup-util utm pull "Windows 11" "/path/in/vm/file.txt" ./local/

  # Push file to VM
  goup-util utm push "Windows 11" ./local/file.txt "/path/in/vm/"`,
}

var utmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all UTM virtual machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUTMCtl("list")
	},
}

var utmStatusCmd = &cobra.Command{
	Use:   "status <vm-name>",
	Short: "Get status of a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUTMCtl("status", args[0])
	},
}

var utmStartCmd = &cobra.Command{
	Use:   "start <vm-name>",
	Short: "Start a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUTMCtl("start", args[0])
	},
}

var utmStopCmd = &cobra.Command{
	Use:   "stop <vm-name>",
	Short: "Stop a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUTMCtl("stop", args[0])
	},
}

var utmIPCmd = &cobra.Command{
	Use:   "ip <vm-name>",
	Short: "Get IP address of a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUTMCtl("ip-address", args[0])
	},
}

var utmExecCmd = &cobra.Command{
	Use:   "exec <vm-name> -- <command> [args...]",
	Short: "Execute a goup-util command in the VM",
	Long: `Execute a goup-util command in the VM.

The command after '--' will be prefixed with 'goup-util' automatically.

Examples:
  # Build for Windows
  goup-util utm exec "Windows 11" -- build windows examples/hybrid-dashboard

  # Generate icons
  goup-util utm exec "Windows 11" -- icons examples/hybrid-dashboard

  # Check config
  goup-util utm exec "Windows 11" -- config`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: goup-util utm exec <vm-name> -- <command> [args...]")
		}

		vmName := args[0]

		// Find the -- separator
		dashIndex := -1
		for i, arg := range args {
			if arg == "--" {
				dashIndex = i
				break
			}
		}

		if dashIndex == -1 || dashIndex == len(args)-1 {
			return fmt.Errorf("missing command after '--'")
		}

		// Build goup-util command
		goupCommand := append([]string{"goup-util"}, args[dashIndex+1:]...)
		cmdStr := strings.Join(goupCommand, " ")

		fmt.Printf("Executing in VM '%s': %s\n\n", vmName, cmdStr)

		// Execute via utmctl
		return runUTMCtl("exec", vmName, "--cmd", cmdStr)
	},
}

var utmTaskCmd = &cobra.Command{
	Use:   "task <vm-name> <task-name>",
	Short: "Execute a Taskfile task in the VM",
	Long: `Execute a Taskfile task in the VM.

This is a convenience wrapper around 'task <taskname>'.

Examples:
  goup-util utm task "Windows 11" build:hybrid:windows
  goup-util utm task "Windows 11" test:all`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vmName := args[0]
		taskName := args[1]

		fmt.Printf("Executing task '%s' in VM '%s'\n\n", taskName, vmName)

		// Execute task via utmctl
		cmdStr := fmt.Sprintf("task %s", taskName)
		return runUTMCtl("exec", vmName, "--cmd", cmdStr)
	},
}

var utmPullCmd = &cobra.Command{
	Use:   "pull <vm-name> <remote-path> <local-path>",
	Short: "Pull a file from the VM to local machine",
	Long: `Pull a file from the VM to local machine.

Examples:
  # Pull MSIX from VM
  goup-util utm pull "Windows 11" "C:\\Users\\User\\goup-util\\examples\\hybrid-dashboard\\.bin\\hybrid-dashboard.msix" ./artifacts/

  # Pull build log
  goup-util utm pull "Windows 11" "/tmp/build.log" ./logs/`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		vmName := args[0]
		remotePath := args[1]
		localPath := args[2]

		fmt.Printf("Pulling from VM '%s':\n", vmName)
		fmt.Printf("  Remote: %s\n", remotePath)
		fmt.Printf("  Local:  %s\n\n", localPath)

		// Use utmctl file pull
		utmctlPath := findUTMCtl()
		pullCmd := exec.Command(utmctlPath, "file", "pull", vmName, remotePath)

		// Create output file
		outFile, err := os.Create(localPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()

		pullCmd.Stdout = outFile
		pullCmd.Stderr = os.Stderr

		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("utmctl file pull failed: %w", err)
		}

		fmt.Printf("✓ File pulled successfully\n")
		return nil
	},
}

var utmPushCmd = &cobra.Command{
	Use:   "push <vm-name> <local-path> <remote-path>",
	Short: "Push a file from local machine to the VM",
	Long: `Push a file from local machine to the VM.

Examples:
  # Push config file
  goup-util utm push "Windows 11" ./config.json "C:\\Users\\User\\config.json"

  # Push test data
  goup-util utm push "Windows 11" ./test-data.zip "/tmp/test-data.zip"`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		vmName := args[0]
		localPath := args[1]
		remotePath := args[2]

		fmt.Printf("Pushing to VM '%s':\n", vmName)
		fmt.Printf("  Local:  %s\n", localPath)
		fmt.Printf("  Remote: %s\n\n", remotePath)

		// Use utmctl file push
		utmctlPath := findUTMCtl()
		pushCmd := exec.Command(utmctlPath, "file", "push", vmName, remotePath)

		// Open input file
		inFile, err := os.Open(localPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer inFile.Close()

		pushCmd.Stdin = inFile
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr

		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("utmctl file push failed: %w", err)
		}

		fmt.Printf("✓ File pushed successfully\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(utmCmd)
	utmCmd.AddCommand(utmListCmd)
	utmCmd.AddCommand(utmStatusCmd)
	utmCmd.AddCommand(utmStartCmd)
	utmCmd.AddCommand(utmStopCmd)
	utmCmd.AddCommand(utmIPCmd)
	utmCmd.AddCommand(utmExecCmd)
	utmCmd.AddCommand(utmTaskCmd)
	utmCmd.AddCommand(utmPullCmd)
	utmCmd.AddCommand(utmPushCmd)
}

// findUTMCtl locates the utmctl binary
func findUTMCtl() string {
	// Try common locations
	locations := []string{
		"/opt/homebrew/bin/utmctl",
		"/usr/local/bin/utmctl",
		"/Applications/UTM.app/Contents/MacOS/utmctl",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// Fallback to PATH
	return "utmctl"
}

// runUTMCtl executes utmctl with the given arguments
func runUTMCtl(args ...string) error {
	utmctlPath := findUTMCtl()

	cmd := exec.Command(utmctlPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
