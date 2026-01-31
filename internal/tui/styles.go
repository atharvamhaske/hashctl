package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Neon magenta/purple theme - flat, no boxes
var (
	// Primary accent - neon magenta/purple
	ColorPrimary = lipgloss.Color("#FF4ECD")
	ColorAccent  = lipgloss.Color("#C77DFF")

	// Secondary colors
	ColorCyan   = lipgloss.Color("#89DDFF")
	ColorGreen  = lipgloss.Color("#C3E88D")
	ColorYellow = lipgloss.Color("#FFCB6B")
	ColorOrange = lipgloss.Color("#F78C6C")
	ColorRed    = lipgloss.Color("#FF5370")

	// Grayscale - muted lavender/gray-purple
	ColorWhite  = lipgloss.Color("#FFFFFF")
	ColorFg     = lipgloss.Color("#E0D4FF")
	ColorMuted  = lipgloss.Color("#8B7FA8")
	ColorDim    = lipgloss.Color("#5C5478")
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
			Foreground(ColorAccent).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorFg)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorDim)
)

// Category/Algorithm styles - flat, no boxes
var (
	// Selected item - bright magenta, no background
	BigSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				PaddingLeft(2)

	// Unselected item - muted, no background
	BigUnselectedStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				PaddingLeft(2)

	// Algorithm list styles - flat
	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	DescStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			PaddingLeft(4)

	CategoryStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	WarningCategoryStyle = lipgloss.NewStyle().
				Foreground(ColorOrange).
				Bold(true)
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
			Foreground(ColorAccent)
)

// Input styles - minimal, no box
var (
	InputStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Bold(true)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)
)

// Result styles - flat
var (
	HashStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	StringStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary)
)

// Help bar
var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(2)
)

// App container - minimal padding
var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

// Cursor for selection
func Cursor() string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("▸ ")
}

// No cursor (spacing)
func NoCursor() string {
	return "  "
}
