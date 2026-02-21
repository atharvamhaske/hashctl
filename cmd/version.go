package cmd

import (
	"fmt"
	"runtime"

	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/atharvamhaske/hashctl/internal/version"
	"github.com/spf13/cobra"
)

var (
	Version   = "v1.2.2"
	BuildDate = "2026-02-21"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run:   runVersion,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Long:  "Check if a newer version of hashctl is available on GitHub.",
	Run:   runCheck,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println()
	label := tui.MutedStyle
	value := tui.ValueStyle

	fmt.Println(label.Render("version   ") + value.Render(Version))
	fmt.Println(label.Render("built     ") + value.Render(BuildDate))
	fmt.Println(label.Render("go        ") + value.Render(runtime.Version()))
	fmt.Println(label.Render("platform  ") + value.Render(runtime.GOOS+"/"+runtime.GOARCH))
	fmt.Println()

	latest, err := version.CheckLatestVersion(Version)
	if err == nil && version.IsUpdateAvailable(Version, latest.TagName) {
		fmt.Println(tui.WarningStyle.Render(
			version.GetUpdateMessage(Version, latest.TagName, latest.URL),
		))
		fmt.Println()
	}
}

func runCheck(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(tui.MutedStyle.Render("Checking for updates..."))
	fmt.Println()

	latest, err := version.CheckLatestVersion(Version)
	if err != nil {
		fmt.Println(tui.ErrorStyle.Render("✗ Failed to check: " + err.Error()))
		fmt.Println()
		return
	}

	if version.IsUpdateAvailable(Version, latest.TagName) {
		fmt.Println(tui.WarningStyle.Render("Update available!"))
		fmt.Println()
		fmt.Println(tui.LabelStyle.Render("Current: ") + tui.ValueStyle.Render(Version))
		fmt.Println(tui.LabelStyle.Render("Latest:  ") + tui.SuccessStyle.Render(latest.TagName))
		fmt.Println()
		fmt.Println(tui.MutedStyle.Render("Download: ") + tui.ValueStyle.Render(latest.URL))
		fmt.Println()
	} else {
		fmt.Println(tui.SuccessStyle.Render("✓ You're on the latest version!"))
		fmt.Println()
		fmt.Println(tui.LabelStyle.Render("Version: ") + tui.ValueStyle.Render(Version))
		fmt.Println()
	}
}
