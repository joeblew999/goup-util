package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/goup-util/pkg/utm"
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

  # List available VMs from gallery
  goup-util utm gallery

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
		return utm.RunUTMCtlInteractive("list")
	},
}

var utmGalleryCmd = &cobra.Command{
	Use:   "gallery",
	Short: "List available VMs from the gallery",
	Long: `List VMs available in the gallery for installation.

The gallery contains pre-configured VM definitions for common operating systems
including Windows 11 ARM, Ubuntu, Debian, and Fedora.

Examples:
  goup-util utm gallery
  goup-util utm gallery --os windows
  goup-util utm gallery --arch arm64`,
	RunE: func(cmd *cobra.Command, args []string) error {
		gallery, err := utm.LoadGallery()
		if err != nil {
			return fmt.Errorf("failed to load gallery: %w", err)
		}

		osFilter, _ := cmd.Flags().GetString("os")
		archFilter, _ := cmd.Flags().GetString("arch")

		vms := gallery.VMs

		// Apply filters
		if osFilter != "" {
			vms = gallery.FilterByOS(osFilter)
		}
		if archFilter != "" {
			filtered := make(map[string]utm.VMEntry)
			for k, v := range vms {
				if v.Arch == archFilter {
					filtered[k] = v
				}
			}
			vms = filtered
		}

		if len(vms) == 0 {
			fmt.Println("No VMs match the filter criteria")
			return nil
		}

		fmt.Println("Available VMs in gallery:")
		fmt.Println()
		for key, vm := range vms {
			fmt.Printf("  %s\n", key)
			fmt.Printf("    Name: %s\n", vm.Name)
			fmt.Printf("    OS:   %s (%s)\n", vm.OS, vm.Arch)
			if vm.Description != "" {
				fmt.Printf("    Desc: %s\n", vm.Description)
			}
			fmt.Printf("    RAM:  %d MB, Disk: %d MB, CPU: %d\n",
				vm.Template.RAM, vm.Template.Disk, vm.Template.CPU)
			if vm.ISO.URL != "" {
				sizeGB := float64(vm.ISO.Size) / 1024 / 1024 / 1024
				fmt.Printf("    ISO:  %.1f GB\n", sizeGB)
			}
			fmt.Println()
		}

		return nil
	},
}

var utmPathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Show UTM paths configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths := utm.GetPaths()
		fmt.Println("UTM Paths:")
		fmt.Printf("  App:   %s\n", paths.App)
		fmt.Printf("  VMs:   %s\n", paths.VMs)
		fmt.Printf("  ISO:   %s\n", paths.ISO)
		fmt.Printf("  Share: %s\n", paths.Share)
		fmt.Println()
		fmt.Printf("utmctl: %s\n", utm.GetUTMCtlPath())
		fmt.Printf("Installed: %v\n", utm.IsUTMInstalled())
		return nil
	},
}

var utmInstallCmd = &cobra.Command{
	Use:   "install [vm-key]",
	Short: "Install UTM app or download VM ISO",
	Long: `Install the UTM application or download a VM ISO from the gallery.

Without arguments, installs the UTM application.
With a VM key, downloads the ISO for that VM.

Examples:
  # Install UTM app
  goup-util utm install

  # Download Windows 11 ISO
  goup-util utm install windows-11-arm

  # Force reinstall UTM
  goup-util utm install --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if len(args) == 0 {
			// Install UTM app
			return utm.InstallUTM(force)
		}

		// Download ISO for specified VM
		return utm.DownloadISO(args[0], force)
	},
}

var utmUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall UTM app",
	RunE: func(cmd *cobra.Command, args []string) error {
		return utm.UninstallUTM()
	},
}

var utmDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check UTM installation status",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := utm.GetInstallStatus()
		if err != nil {
			return err
		}

		fmt.Println("UTM Installation Status:")
		fmt.Println()

		if status.Installed {
			fmt.Printf("  ✓ UTM installed at %s\n", status.InstalledPath)
			if status.InstalledVersion != "" {
				fmt.Printf("    Version: %s\n", status.InstalledVersion)
			}
		} else {
			fmt.Printf("  ✗ UTM not installed\n")
			fmt.Printf("    Run: goup-util utm install\n")
		}

		fmt.Printf("  Gallery version: %s\n", status.GalleryVersion)

		if status.UpdateAvailable {
			fmt.Printf("  ⚠ Update available: %s\n", status.GalleryVersion)
			fmt.Printf("    Run: goup-util utm install --force\n")
		}

		// Check directories
		paths := utm.GetPaths()
		fmt.Println()
		fmt.Println("Directories:")
		checkDir("VMs", paths.VMs)
		checkDir("ISO", paths.ISO)
		checkDir("Share", paths.Share)

		return nil
	},
}

func checkDir(name, path string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("  ✓ %s: %s\n", name, path)
	} else {
		fmt.Printf("  ✗ %s: %s (missing)\n", name, path)
	}
}

var utmStatusCmd = &cobra.Command{
	Use:   "status <vm-name>",
	Short: "Get status of a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := utm.GetVMStatus(args[0])
		if err != nil {
			return err
		}
		fmt.Println(status)
		return nil
	},
}

var utmStartCmd = &cobra.Command{
	Use:   "start <vm-name>",
	Short: "Start a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return utm.StartVM(args[0])
	},
}

var utmStopCmd = &cobra.Command{
	Use:   "stop <vm-name>",
	Short: "Stop a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return utm.StopVM(args[0])
	},
}

var utmIPCmd = &cobra.Command{
	Use:   "ip <vm-name>",
	Short: "Get IP address of a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip, err := utm.GetVMIP(args[0])
		if err != nil {
			return err
		}
		fmt.Println(ip)
		return nil
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

		return utm.ExecInVM(vmName, cmdStr)
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

		cmdStr := fmt.Sprintf("task %s", taskName)
		return utm.ExecInVM(vmName, cmdStr)
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

		if err := utm.PullFile(vmName, remotePath, localPath); err != nil {
			return err
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

		if err := utm.PushFile(vmName, localPath, remotePath); err != nil {
			return err
		}

		fmt.Printf("✓ File pushed successfully\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(utmCmd)
	utmCmd.AddCommand(utmListCmd)
	utmCmd.AddCommand(utmGalleryCmd)
	utmCmd.AddCommand(utmPathsCmd)
	utmCmd.AddCommand(utmInstallCmd)
	utmCmd.AddCommand(utmUninstallCmd)
	utmCmd.AddCommand(utmDoctorCmd)
	utmCmd.AddCommand(utmStatusCmd)
	utmCmd.AddCommand(utmStartCmd)
	utmCmd.AddCommand(utmStopCmd)
	utmCmd.AddCommand(utmIPCmd)
	utmCmd.AddCommand(utmExecCmd)
	utmCmd.AddCommand(utmTaskCmd)
	utmCmd.AddCommand(utmPullCmd)
	utmCmd.AddCommand(utmPushCmd)

	// Gallery filters
	utmGalleryCmd.Flags().String("os", "", "Filter by OS (windows, linux)")
	utmGalleryCmd.Flags().String("arch", "", "Filter by architecture (arm64, amd64)")

	// Install flags
	utmInstallCmd.Flags().Bool("force", false, "Force reinstall/redownload")
}
