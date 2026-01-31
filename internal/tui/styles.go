package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Clean, minimal color palette inspired by charmbracelet projects
var (
	// Primary accent - soft pink/magenta
	ColorPrimary = lipgloss.Color("#FF6B9D")
	ColorAccent  = lipgloss.Color("#C792EA")

	// Secondary colors
	ColorCyan   = lipgloss.Color("#89DDFF")
	ColorGreen  = lipgloss.Color("#C3E88D")
	ColorYellow = lipgloss.Color("#FFCB6B")
	ColorOrange = lipgloss.Color("#F78C6C")
	ColorRed    = lipgloss.Color("#FF5370")

	// Grayscale
	ColorWhite  = lipgloss.Color("#FFFFFF")
	ColorFg     = lipgloss.Color("#EEFFFF")
	ColorMuted  = lipgloss.Color("#676E95")
	ColorDim    = lipgloss.Color("#4A5568")
	ColorBorder = lipgloss.Color("#3B4252")
	ColorBg     = lipgloss.Color("#0F111A")
	ColorBgAlt  = lipgloss.Color("#1A1C25")
)

// Logo and branding
var (
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	LogoAccent = lipgloss.NewStyle().
			Foreground(ColorAccent)
)

// Text styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorFg)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorDim)
)

// List styles - Big text for categories
var (
	// Big highlighted selected item - smaller box, bigger text
	BigSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Padding(1, 3).
				Background(ColorBgAlt).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Width(45).
				Height(3)

	// Big unselected item - smaller box
	BigUnselectedStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Padding(1, 3).
				Width(45).
				Height(3)

	// Algorithm list styles - bigger text
	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(1, 3).
			Background(ColorBgAlt)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Padding(1, 3)

	DescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true).
			PaddingLeft(4)

	CategoryStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			MarginTop(1)

	WarningCategoryStyle = lipgloss.NewStyle().
				Foreground(ColorOrange).
				Bold(true).
				MarginTop(1)
)

// Status styles
var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(ColorCyan)
)

// Box styles
var (
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	ResultBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGreen).
			Padding(1, 2)

	InputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorCyan).
			Padding(2, 3).
			Width(60)

	// Big input text style (no placeholder, big text)
	BigInputStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Bold(true).
			Width(60).
			Height(3)
)

// Help bar
var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(2)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ColorDim)
)

// Hash display
var (
	HashStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(ColorCyan)

	StringStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)
)

// App container
var (
	AppStyle = lipgloss.NewStyle().
		Padding(1, 2)
)

// Cursor for selection
func Cursor() string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render("▸ ")
}

// No cursor (spacing)
func NoCursor() string {
	return "  "
}

// Badge renders a small badge
func Badge(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render("[" + text + "]")
}

// Divider renders a horizontal line
func Divider(width int) string {
	line := ""
	for i := 0; i < width; i++ {
		line += "─"
	}
	return DimStyle.Render(line)
}
