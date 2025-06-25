package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goup-util",
	Short: "A CLI tool for managing Android and iOS SDKs",
	Long:  "Go Up Util is a CLI tool for managing Android and iOS SDK installations using JSON files.",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
