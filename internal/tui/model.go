package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atharvamhaske/hashctl/internal/hasher"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// State represents the current TUI screen
type State int

const (
	StateAlgorithmSelect State = iota
	StateInputMode
	StateTextInput
	StateHashing
	StateResults
)

// InputMode represents what we're hashing
type InputMode int

const (
	InputModeString InputMode = iota
	InputModeFile
)

// Model is the main TUI model
type Model struct {
	state     State
	inputMode InputMode

	// Algorithm selection
	algorithms     []hasher.Algorithm
	algorithmIndex int
	selectedAlgo   hasher.Algorithm

	// Input
	textInput textinput.Model
	files     []string

	// Hashing
	spinner   spinner.Model
	isHashing bool
	hashStart time.Time

	// Results
	results []hasher.Result

	// UI dimensions
	width  int
	height int
	err    error

	// Options
	opts hasher.Options
}

// Messages
type hashCompleteMsg struct {
	results []hasher.Result
}

type hashErrorMsg struct {
	err error
}

// NewModel creates a new TUI model
func NewModel() Model {
	algs := hasher.GetSortedAlgorithms()

	ti := textinput.New()
	ti.Placeholder = "type here..."
	ti.CharLimit = 10000
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = SpinnerStyle

	return Model{
		state:          StateAlgorithmSelect,
		algorithms:     algs,
		algorithmIndex: 0,
		textInput:      ti,
		spinner:        s,
		opts:           hasher.DefaultOptions(),
		width:          80,
		height:         24,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles all messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case StateAlgorithmSelect:
			return m.handleAlgorithmSelect(msg)
		case StateInputMode:
			return m.handleInputMode(msg)
		case StateTextInput:
			return m.handleTextInput(msg)
		case StateHashing:
			// No input during hashing
			return m, nil
		case StateResults:
			return m.handleResults(msg)
		}

	case spinner.TickMsg:
		if m.isHashing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case hashCompleteMsg:
		m.isHashing = false
		m.results = msg.results
		m.state = StateResults
		return m, nil

	case hashErrorMsg:
		m.isHashing = false
		m.err = msg.err
		m.state = StateResults
		return m, nil
	}

	return m, nil
}

func (m Model) handleAlgorithmSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.algorithmIndex > 0 {
			m.algorithmIndex--
		}
	case "down", "j":
		if m.algorithmIndex < len(m.algorithms)-1 {
			m.algorithmIndex++
		}
	case "enter", " ":
		m.selectedAlgo = m.algorithms[m.algorithmIndex]
		m.opts.Algorithm = getAlgorithmKey(m.selectedAlgo.Name)
		m.state = StateInputMode
	case "home", "g":
		m.algorithmIndex = 0
	case "end", "G":
		m.algorithmIndex = len(m.algorithms) - 1
	}
	return m, nil
}

func (m Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc":
		m.state = StateAlgorithmSelect
	case "s", "1":
		m.inputMode = InputModeString
		m.textInput.Placeholder = "enter text to hash..."
		m.textInput.Reset()
		m.textInput.Focus()
		m.state = StateTextInput
		return m, textinput.Blink
	case "f", "2":
		m.inputMode = InputModeFile
		m.textInput.Placeholder = "enter file path..."
		m.textInput.Reset()
		m.textInput.Focus()
		m.state = StateTextInput
		return m, textinput.Blink
	}
	return m, nil
}

func (m Model) handleTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.textInput.Reset()
		m.state = StateInputMode
		return m, nil
	case "enter":
		input := strings.TrimSpace(m.textInput.Value())
		if input == "" {
			return m, nil
		}
		m.state = StateHashing
		m.isHashing = true
		m.hashStart = time.Now()

		if m.inputMode == InputModeString {
			return m, tea.Batch(m.spinner.Tick, m.doHashString(input))
		}
		m.files = []string{input}
		return m, tea.Batch(m.spinner.Tick, m.doHashFiles())
	default:
		// IMPORTANT: Pass all other keys to the text input!
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc", "r":
		m.state = StateAlgorithmSelect
		m.results = nil
		m.err = nil
		m.textInput.Reset()
		m.files = nil
	case "n":
		// New hash with same algorithm
		m.results = nil
		m.err = nil
		m.textInput.Reset()
		m.state = StateInputMode
	}
	return m, nil
}

func (m Model) doHashString(input string) tea.Cmd {
	return func() tea.Msg {
		result := hasher.HashString(input, m.opts)
		return hashCompleteMsg{results: []hasher.Result{result}}
	}
}

func (m Model) doHashFiles() tea.Cmd {
	return func() tea.Msg {
		var results []hasher.Result
		hasher.HashFiles(m.files, m.opts, func(r hasher.Result) {
			results = append(results, r)
		})
		return hashCompleteMsg{results: results}
	}
}

// View renders the TUI
func (m Model) View() string {
	var s strings.Builder

	switch m.state {
	case StateAlgorithmSelect:
		s.WriteString(m.viewAlgorithmSelect())
	case StateInputMode:
		s.WriteString(m.viewInputMode())
	case StateTextInput:
		s.WriteString(m.viewTextInput())
	case StateHashing:
		s.WriteString(m.viewHashing())
	case StateResults:
		s.WriteString(m.viewResults())
	}

	return AppStyle.Render(s.String())
}

func (m Model) viewAlgorithmSelect() string {
	var s strings.Builder

	// Header
	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(LogoAccent.Render(" ⟡"))
	s.WriteString("\n")
	s.WriteString(SubtitleStyle.Render("compute cryptographic hashes for strings & files"))
	s.WriteString("\n\n")

	// Group by category
	byCategory := hasher.GetAlgorithmsByCategory()
	categories := []hasher.Category{
		hasher.CategoryChecksum,
		hasher.CategoryFastHash,
		hasher.CategoryPasswordHash,
	}

	globalIndex := 0
	for _, cat := range categories {
		algs, ok := byCategory[cat]
		if !ok || len(algs) == 0 {
			continue
		}

		// Category header
		catStyle := CategoryStyle
		if cat == hasher.CategoryPasswordHash {
			catStyle = WarningCategoryStyle
		}
		s.WriteString(catStyle.Render(cat.String()))
		s.WriteString("\n")

		for _, alg := range algs {
			isSelected := m.algorithmIndex == globalIndex

			if isSelected {
				s.WriteString(Cursor())
				s.WriteString(SelectedStyle.Render(alg.Name))
				s.WriteString("\n")
				s.WriteString("   ")
				s.WriteString(DescStyle.Render(alg.Description))
				if alg.IsPasswordHash {
					s.WriteString("\n   ")
					s.WriteString(WarningStyle.Render("⚠ slow on large inputs"))
				}
			} else {
				s.WriteString(NoCursor())
				s.WriteString(UnselectedStyle.Render(alg.Name))
			}
			s.WriteString("\n")
			globalIndex++
		}
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/↓ select • enter confirm • q quit"))

	return s.String()
}

func (m Model) viewInputMode() string {
	var s strings.Builder

	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(LogoAccent.Render(" ⟡ "))
	s.WriteString(LabelStyle.Render(m.selectedAlgo.Name))
	s.WriteString("\n\n")

	s.WriteString(SubtitleStyle.Render("what do you want to hash?"))
	s.WriteString("\n\n")

	s.WriteString("  ")
	s.WriteString(Badge("s", ColorCyan))
	s.WriteString("  ")
	s.WriteString(ValueStyle.Render("hash a string"))
	s.WriteString("\n\n")

	s.WriteString("  ")
	s.WriteString(Badge("f", ColorCyan))
	s.WriteString("  ")
	s.WriteString(ValueStyle.Render("hash a file"))
	s.WriteString("\n\n")

	s.WriteString(HelpStyle.Render("s string • f file • esc back • q quit"))

	return s.String()
}

func (m Model) viewTextInput() string {
	var s strings.Builder

	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(LogoAccent.Render(" ⟡ "))
	s.WriteString(LabelStyle.Render(m.selectedAlgo.Name))
	s.WriteString("\n\n")

	var label string
	if m.inputMode == InputModeString {
		label = "enter text:"
	} else {
		label = "enter file path:"
	}
	s.WriteString(SubtitleStyle.Render(label))
	s.WriteString("\n\n")

	// Input field
	inputBox := InputBoxStyle.Render(m.textInput.View())
	s.WriteString(inputBox)
	s.WriteString("\n\n")

	s.WriteString(HelpStyle.Render("enter hash • esc back"))

	return s.String()
}

func (m Model) viewHashing() string {
	var s strings.Builder

	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString("\n\n")

	s.WriteString(m.spinner.View())
	s.WriteString(" ")
	s.WriteString(MutedStyle.Render("computing hash..."))
	s.WriteString("\n\n")

	elapsed := time.Since(m.hashStart)
	s.WriteString(DimStyle.Render(fmt.Sprintf("elapsed: %s", elapsed.Round(time.Millisecond))))

	return s.String()
}

func (m Model) viewResults() string {
	var s strings.Builder

	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(" ")
	s.WriteString(SuccessStyle.Render("✓"))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(ErrorStyle.Render("error: " + m.err.Error()))
		s.WriteString("\n")
	} else {
		s.WriteString(LabelStyle.Render(m.selectedAlgo.Name))
		s.WriteString("\n")
		s.WriteString(Divider(50))
		s.WriteString("\n\n")

		for _, r := range m.results {
			if r.Error != nil {
				s.WriteString(ErrorStyle.Render("✗ " + r.Input))
				s.WriteString("\n")
				s.WriteString(MutedStyle.Render("  " + r.Error.Error()))
				s.WriteString("\n\n")
			} else {
				// Input label
				if r.IsFile {
					s.WriteString(FileStyle.Render("file: "))
					s.WriteString(ValueStyle.Render(r.Input))
				} else {
					s.WriteString(StringStyle.Render("text: "))
					s.WriteString(MutedStyle.Render("\"" + truncate(r.Input, 30) + "\""))
				}
				s.WriteString("\n\n")

				// Hash result in a box
				hashBox := ResultBoxStyle.Render(HashStyle.Render(r.Hash))
				s.WriteString(hashBox)
				s.WriteString("\n\n")

				s.WriteString(DimStyle.Render(fmt.Sprintf("computed in %s", r.Duration.Round(time.Microsecond))))
				s.WriteString("\n")
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("n new hash • r restart • q quit"))

	return s.String()
}

// Helpers
func getAlgorithmKey(name string) string {
	for key, alg := range hasher.Registry {
		if alg.Name == name {
			return key
		}
	}
	return strings.ToLower(name)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Run starts the TUI application
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
