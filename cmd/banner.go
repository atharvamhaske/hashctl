package cmd

import (
	"fmt"

	"github.com/atharvamhaske/hashctl/internal/tui"
)

// GenerateBanner builds the PMG-style 2-line banner string.
//
//	█▄█ ▄▀▄ █▀▀ █▄█ █▀▀ ▀█▀ █    From Atharva Mhaske (github.com/atharvamhaske/hashctl)
//	█ █ █▀█ ▀▀█ █ █ █▄▄  █  █▄▄  version: v1.2.0
func GenerateBanner() string {
	line1 := fmt.Sprintf("%s\tFrom %s %s",
		tui.BrandStyle.Render(tui.BrandLine1),
		tui.LabelStyle.Render("Atharva Mhaske"),
		tui.MutedStyle.Render("(github.com/atharvamhaske/hashctl)"),
	)

	line2 := fmt.Sprintf("%s\t%s: %s",
		tui.BrandStyle.Render(tui.BrandLine2),
		tui.MutedStyle.Render("version"),
		tui.ValueStyle.Render(Version),
	)

	return "\n" + line1 + "\n" + line2 + "\n"
}

// PrintBanner prints the 2-line banner to stdout.
func PrintBanner() {
	fmt.Print(GenerateBanner())
}

// PrintUsage prints the banner followed by compact usage info.
// Registered as the root command's help function.
func PrintUsage() {
	PrintBanner()

	fmt.Println()
	fmt.Println(tui.MutedStyle.Render("Usage:"))
	fmt.Println(tui.ValueStyle.Render("  hashctl") + tui.MutedStyle.Render("          launch interactive TUI"))
	fmt.Println(tui.ValueStyle.Render("  hashctl list") + tui.MutedStyle.Render("      list supported algorithms"))
	fmt.Println(tui.ValueStyle.Render("  hashctl version") + tui.MutedStyle.Render("   show version info"))
	fmt.Println(tui.ValueStyle.Render("  hashctl check") + tui.MutedStyle.Render("     check for updates"))
	fmt.Println()
	fmt.Println(tui.MutedStyle.Render("Flags:"))
	fmt.Println(tui.MutedStyle.Render("  -h, --help   show this message"))
	fmt.Println()
}
