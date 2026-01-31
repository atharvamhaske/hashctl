// Package cmd provides the CLI entry point for hashctl
package cmd

import (
	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hashctl",
	Short: "Interactive hashing TUI",
	Long: `hashctl is an interactive terminal UI for computing cryptographic hashes.

Launch hashctl to select from 20+ algorithms and hash strings or files
with a beautiful, keyboard-driven interface.

Supported algorithms include SHA-256, SHA-512, BLAKE2, SHA-3, MD5, 
bcrypt, Argon2id, and more.`,
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
}
