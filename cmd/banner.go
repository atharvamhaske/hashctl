package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// 4-col × 5-row pixel font – consistent 1px strokes, single-char blocks.
// Total: 7 letters × 4 + 6 gaps × 1 = 34 chars wide, 5 rows tall.
var font = map[byte][5]string{
	'H': {"█  █", "█  █", "████", "█  █", "█  █"},
	'A': {" ██ ", "█  █", "████", "█  █", "█  █"},
	'S': {"████", "█   ", "████", "   █", "████"},
	'C': {"████", "█   ", "█   ", "█   ", "████"},
	'T': {"████", " █  ", " █  ", " █  ", " █  "},
	'L': {"█   ", "█   ", "█   ", "█   ", "████"},
}

// renderArt builds the 5-line HASHCTL pixel art with 1-char gaps.
func renderArt() string {
	const word = "HASHCTL"
	var rows [5]string
	for r := 0; r < 5; r++ {
		for i := 0; i < len(word); i++ {
			if i > 0 {
				rows[r] += " "
			}
			rows[r] += font[word[i]][r]
		}
	}
	return strings.Join(rows[:], "\n")
}

// ── Styles ────────────────────────────────────────────────────────────────────

var bs = struct {
	art   lipgloss.Style
	title lipgloss.Style
	muted lipgloss.Style
	value lipgloss.Style
}{
	art:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4ECD")).Bold(true),
	title: lipgloss.NewStyle().Foreground(lipgloss.Color("#C77DFF")).Bold(true),
	muted: lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7FA8")),
	value: lipgloss.NewStyle().Foreground(lipgloss.Color("#E0D4FF")),
}

// ── Public functions ──────────────────────────────────────────────────────────

// PrintBanner prints the compact HASHCTL pixel-art banner with
// author + version info to the right (PMG-style, fits in 80 cols).
func PrintBanner() {
	left := bs.art.Render(renderArt())

	info := strings.Join([]string{
		"",
		bs.title.Render("From Atharva Mhaske"),
		bs.muted.Render("github.com/atharvamhaske/hashctl"),
		bs.muted.Render("version: ") + bs.value.Render(Version),
		"",
	}, "\n")

	right := lipgloss.NewStyle().PaddingLeft(3).Render(info)

	fmt.Println()
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, left, right))
	fmt.Println()
}

// PrintUsage prints the banner followed by usage info.
// Used as the root command's custom help function.
func PrintUsage() {
	PrintBanner()

	fmt.Println(bs.muted.Render("Usage:"))
	fmt.Println(bs.value.Render("  hashctl") +
		bs.muted.Render("          launch interactive TUI"))
	fmt.Println(bs.value.Render("  hashctl list") +
		bs.muted.Render("      list supported algorithms"))
	fmt.Println(bs.value.Render("  hashctl version") +
		bs.muted.Render("   show version info"))
	fmt.Println(bs.value.Render("  hashctl check") +
		bs.muted.Render("     check for updates"))
	fmt.Println()
	fmt.Println(bs.muted.Render("Flags:"))
	fmt.Println(bs.muted.Render("  -h, --help   show this message"))
	fmt.Println()
}
