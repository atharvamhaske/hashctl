// Package cmd provides the CLI entry point for hashctl
package cmd

import (
	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "hashctl",
	Short:         "Interactive hashing TUI",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(checkCmd)

	// Replace Cobra's default help with our PMG-style banner + usage
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		PrintUsage()
	})
}
