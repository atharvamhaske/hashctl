package cmd

// Pixel-art bitmaps (each pixel = "▓▓", 4 cols × 5 rows per letter, 2-char gap between letters)
// Total banner width = 7 × 8 + 6 × 2 = 68 chars per row.
//
//   H       A       S       H       C       T       L
//   1001    0110    1110    1001    1110    1111    1000
//   1001    1001    1000    1001    1000    0110    1000
//   1111    1111    1110    1111    1000    0110    1000
//   1001    1001    0001    1001    1000    0110    1000
//   1001    1001    0111    1001    1110    0110    1111

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII art rows for HASHCTL (pre-rendered, each row is exactly 68 chars before trim)
var asciiRows = [5]string{
	"▓▓    ▓▓    ▓▓▓▓    ▓▓▓▓▓▓    ▓▓    ▓▓  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▓  ▓▓      ",
	"▓▓    ▓▓  ▓▓    ▓▓  ▓▓        ▓▓    ▓▓  ▓▓          ▓▓▓▓    ▓▓      ",
	"▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▓  ▓▓          ▓▓▓▓    ▓▓      ",
	"▓▓    ▓▓  ▓▓    ▓▓        ▓▓  ▓▓    ▓▓  ▓▓          ▓▓▓▓    ▓▓      ",
	"▓▓    ▓▓  ▓▓    ▓▓    ▓▓▓▓▓▓  ▓▓    ▓▓  ▓▓▓▓▓▓      ▓▓▓▓    ▓▓▓▓▓▓▓▓",
}

// bannerStyles defines reusable lipgloss styles for the banner.
var bannerStyles = struct {
	art     lipgloss.Style
	title   lipgloss.Style
	muted   lipgloss.Style
	value   lipgloss.Style
	success lipgloss.Style
}{
	art:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4ECD")).Bold(true),
	title:   lipgloss.NewStyle().Foreground(lipgloss.Color("#C77DFF")).Bold(true),
	muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7FA8")),
	value:   lipgloss.NewStyle().Foreground(lipgloss.Color("#E0D4FF")),
	success: lipgloss.NewStyle().Foreground(lipgloss.Color("#C3E88D")).Bold(true),
}

// PrintBanner prints the PMG-style startup banner:
//
//	[HASHCTL ASCII ART]   From Atharva Mhaske
//	                      github.com/atharvamhaske/hashctl
//	                      version: v1.2.0 · go1.21 · linux/amd64
//
// It is used by `hashctl version` and the root-command help template.
func PrintBanner() {
	s := bannerStyles

	// ── left block: coloured ASCII art ────────────────────────────────────────
	var leftLines [5]string
	for i, row := range asciiRows {
		leftLines[i] = s.art.Render(strings.TrimRight(row, " "))
	}
	leftBlock := strings.Join(leftLines[:], "\n")

	// ── right block: author / repo / version metadata ─────────────────────────
	infoLines := []string{
		"",
		s.title.Render("From Atharva Mhaske") +
			s.muted.Render("  (github.com/atharvamhaske/hashctl)"),
		s.muted.Render("version: ") + s.value.Render(Version) +
			s.muted.Render("  built: ") + s.value.Render(BuildDate),
		s.muted.Render("go:      ") + s.value.Render(runtime.Version()) +
			s.muted.Render("  ") + s.value.Render(runtime.GOOS+"/"+runtime.GOARCH),
		"",
	}
	rightBlock := lipgloss.NewStyle().PaddingLeft(4).
		Render(strings.Join(infoLines, "\n"))

	// ── join side-by-side ─────────────────────────────────────────────────────
	banner := lipgloss.JoinHorizontal(lipgloss.Center, leftBlock, rightBlock)

	fmt.Println()
	fmt.Println(banner)
	fmt.Println()
}

// PrintUsage prints the banner followed by Cobra-style usage info.
// Used as the root command's custom help function.
func PrintUsage() {
	PrintBanner()

	s := bannerStyles
	fmt.Println(s.muted.Render("Usage:"))
	fmt.Println(s.value.Render("  hashctl") +
		s.muted.Render("          launch the interactive TUI"))
	fmt.Println(s.value.Render("  hashctl list") +
		s.muted.Render("      list all supported algorithms"))
	fmt.Println(s.value.Render("  hashctl version") +
		s.muted.Render("   show version and check for updates"))
	fmt.Println(s.value.Render("  hashctl check") +
		s.muted.Render("     check for a newer release"))
	fmt.Println()
	fmt.Println(s.muted.Render("Flags:"))
	fmt.Println(s.muted.Render("  -h, --help   show this message"))
	fmt.Println()
	fmt.Println(s.muted.Render("Source · https://github.com/atharvamhaske/hashctl"))
	fmt.Println()
}

