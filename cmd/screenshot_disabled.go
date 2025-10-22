//go:build !screenshot
// +build !screenshot

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Screenshot support not available in this build",
	Long: `Screenshot support requires CGO and is not available in this build.

To build with screenshot support:
  CGO_ENABLED=1 go build -tags screenshot .

Or use the task:
  task build:with-screenshot`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("screenshot support not available in this build\nRebuild with: CGO_ENABLED=1 go build -tags screenshot")
	},
}

func init() {
	rootCmd.AddCommand(screenshotCmd)
}
