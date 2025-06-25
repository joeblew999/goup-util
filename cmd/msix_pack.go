package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var msixPackCmd = &cobra.Command{
	Use:   "msix-pack",
	Short: "Package an MSIX file (Windows only)",
	Long:  "Package an MSIX file using the msix toolkit. This command only works on Windows.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if we're on Windows
		if runtime.GOOS != "windows" {
			fmt.Println("Skipping MSIX packaging: This command only works on Windows.")
			return nil
		}

		// Get flags
		directory, _ := cmd.Flags().GetString("directory")
		packagePath, _ := cmd.Flags().GetString("package")

		if directory == "" || packagePath == "" {
			return fmt.Errorf("both --directory and --package flags are required")
		}

		// Check if msix command is available
		if _, err := exec.LookPath("msix"); err != nil {
			return fmt.Errorf("msix command not found. Please install the MSIX toolkit")
		}

		// Run the msix pack command
		msixCmd := exec.Command("msix", "pack", "-d", directory, "-p", packagePath)
		msixCmd.Stdout = os.Stdout
		msixCmd.Stderr = os.Stderr

		fmt.Printf("Packaging MSIX: %s\n", packagePath)
		if err := msixCmd.Run(); err != nil {
			return fmt.Errorf("failed to package MSIX: %w", err)
		}

		fmt.Printf("Successfully created MSIX package: %s\n", packagePath)
		return nil
	},
}

func init() {
	msixPackCmd.Flags().StringP("directory", "d", "", "Directory containing the app files")
	msixPackCmd.Flags().StringP("package", "p", "", "Output path for the MSIX package")
	msixPackCmd.MarkFlagRequired("directory")
	msixPackCmd.MarkFlagRequired("package")
	rootCmd.AddCommand(msixPackCmd)
}
