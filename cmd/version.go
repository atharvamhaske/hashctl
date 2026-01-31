package cmd

import (
	"fmt"
	"runtime"

	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/spf13/cobra"
)

var (
	Version   = "v0.2.0"
	BuildDate = "2026-01-31"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(tui.LogoStyle.Render("hashctl") + tui.LogoAccent.Render(" ⟡"))
	fmt.Println()

	label := tui.MutedStyle
	value := tui.ValueStyle

	fmt.Println(label.Render("version   ") + value.Render(Version))
	fmt.Println(label.Render("built     ") + value.Render(BuildDate))
	fmt.Println(label.Render("go        ") + value.Render(runtime.Version()))
	fmt.Println(label.Render("platform  ") + value.Render(runtime.GOOS+"/"+runtime.GOARCH))
	fmt.Println()
}
