// Package cmd provides the CLI entry point for hashctl
package cmd

import (
	"fmt"

	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "hashctl",
	Short:         "Interactive hashing TUI",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show banner before launching TUI
		PrintBanner()
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

	// Custom help without banner (banner only shows on `hashctl` command)
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		dim := tui.MutedStyle
		val := tui.ValueStyle

		fmt.Println()
		fmt.Println(dim.Render("Usage:"))
		fmt.Println(val.Render("  hashctl") + dim.Render("          launch interactive TUI"))
		fmt.Println(val.Render("  hashctl list") + dim.Render("      list supported algorithms"))
		fmt.Println(val.Render("  hashctl version") + dim.Render("   show version info"))
		fmt.Println(val.Render("  hashctl check") + dim.Render("     check for updates"))
		fmt.Println()
		fmt.Println(dim.Render("Flags:"))
		fmt.Println(dim.Render("  -h, --help   show this message"))
		fmt.Println()
	})
}
