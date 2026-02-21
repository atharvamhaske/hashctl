package cmd

import (
	"fmt"

	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/atharvamhaske/hashctl/internal/version"
	"github.com/spf13/cobra"
)

var (
	Version   = "v1.2.0"
	BuildDate = "2026-02-21"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and check for updates",
	Run:   runVersion,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Long:  "Check if a newer version of hashctl is available on GitHub.",
	Run:   runCheck,
}

func runVersion(cmd *cobra.Command, args []string) {
	// Full PMG-style banner with version metadata embedded
	PrintBanner()

	// Update check (non-blocking best-effort)
	latest, err := version.CheckLatestVersion(Version)
	if err == nil && version.IsUpdateAvailable(Version, latest.TagName) {
		fmt.Println(tui.WarningStyle.Render(
			version.GetUpdateMessage(Version, latest.TagName, latest.URL),
		))
		fmt.Println()
	}
}

func runCheck(cmd *cobra.Command, args []string) {
	PrintBanner()

	fmt.Println(tui.MutedStyle.Render("Checking for updates..."))
	fmt.Println()

	latest, err := version.CheckLatestVersion(Version)
	if err != nil {
		fmt.Println(tui.ErrorStyle.Render("✗ Failed to check for updates: " + err.Error()))
		fmt.Println()
		return
	}

	if version.IsUpdateAvailable(Version, latest.TagName) {
		fmt.Println(tui.WarningStyle.Render("Update available!"))
		fmt.Println()
		fmt.Println(tui.LabelStyle.Render("Current version:") + " " + tui.ValueStyle.Render(Version))
		fmt.Println(tui.LabelStyle.Render("Latest version: ") + " " + tui.SuccessStyle.Render(latest.TagName))
		fmt.Println()
		fmt.Println(tui.MutedStyle.Render("Download: ") + tui.ValueStyle.Render(latest.URL))
		fmt.Println()
	} else {
		fmt.Println(tui.SuccessStyle.Render("✓ You're on the latest version!"))
		fmt.Println()
		fmt.Println(tui.LabelStyle.Render("Version: ") + " " + tui.ValueStyle.Render(Version))
		fmt.Println()
	}
}
