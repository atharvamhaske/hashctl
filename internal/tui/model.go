package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atharvamhaske/hashctl/internal/hasher"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ── State machine ────────────────────────────────────────────────────────────

// State represents the current TUI screen.
type State int

const (
	StateCategorySelect State = iota
	StateAlgorithmSelect
	StateInputMode
	StateTextInput
	StateHashing
	StateResults
)

// InputMode represents what we're hashing.
type InputMode int

const (
	InputModeString InputMode = iota
	InputModeFile
)

// ── list.Item implementations ────────────────────────────────────────────────

// categoryItem is shown on the category-select screen.
type categoryItem struct {
	category    hasher.Category
	title       string
	description string
}

func (i categoryItem) Title() string       { return i.title }
func (i categoryItem) Description() string { return i.description }
func (i categoryItem) FilterValue() string { return i.title }

// algorithmItem is shown on the algorithm-select screen.
type algorithmItem struct {
	alg hasher.Algorithm
}

func (i algorithmItem) Title() string {
	return strings.ToUpper(i.alg.Name)
}

func (i algorithmItem) Description() string {
	if i.alg.IsPasswordHash {
		return i.alg.Description + "  ⚠ slow by design"
	}
	return i.alg.Description
}

func (i algorithmItem) FilterValue() string { return i.alg.Name }

// ── Model ────────────────────────────────────────────────────────────────────

// Model is the main Bubble Tea model.
type Model struct {
	state     State
	inputMode InputMode

	// Fancy list components
	categoryList  list.Model
	algorithmList list.Model

	// Resolved selections
	selectedCategory hasher.Category
	selectedAlgo     hasher.Algorithm

	// Text input (string / file path)
	textInput textinput.Model
	files     []string

	// Hashing progress
	spinner   spinner.Model
	isHashing bool
	hashStart time.Time

	// Results
	results []hasher.Result
	err     error

	// Terminal dimensions
	width  int
	height int

	// Hashing options
	opts hasher.Options
}

// ── Async messages ───────────────────────────────────────────────────────────

type hashCompleteMsg struct{ results []hasher.Result }
type hashErrorMsg struct{ err error }

// ── List constructors ────────────────────────────────────────────────────────

func makeCategoryList() list.Model {
	items := []list.Item{
		categoryItem{
			category:    hasher.CategoryChecksum,
			title:       "CHECKSUMS",
			description: "Non-cryptographic checksums — CRC32 and variants",
		},
		categoryItem{
			category:    hasher.CategoryFastHash,
			title:       "FAST CRYPTOGRAPHIC HASHES",
			description: "MD5 · SHA-1 · SHA-2 · SHA-3 · BLAKE2 · RIPEMD-160",
		},
		categoryItem{
			category:    hasher.CategoryPasswordHash,
			title:       "PASSWORD HASHING / KDFS",
			description: "bcrypt · Argon2id — slow by design, resistant to GPU attacks",
		},
	}

	d := NewListDelegate(true)
	l := list.New(items, d, 70, 10)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	return l
}

func makeAlgorithmList(category hasher.Category) list.Model {
	all := hasher.GetSortedAlgorithms()
	var items []list.Item
	for _, alg := range all {
		if alg.Category == category {
			items = append(items, algorithmItem{alg: alg})
		}
	}

	d := NewListDelegate(true)
	l := list.New(items, d, 70, 18)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(true) // shows "X items" count
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()

	// Style the status bar and filter prompt with our theme
	l.Styles.StatusBar = MutedStyle
	l.Styles.FilterPrompt = LabelStyle
	l.Styles.FilterCursor = SelectedStyle
	return l
}

// ── NewModel ─────────────────────────────────────────────────────────────────

// NewModel creates the initial TUI model.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 10000
	ti.Width = 68

	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = SpinnerStyle

	return Model{
		state:         StateCategorySelect,
		categoryList:  makeCategoryList(),
		algorithmList: list.New([]list.Item{}, NewListDelegate(true), 70, 18),
		textInput:     ti,
		spinner:       s,
		opts:          hasher.DefaultOptions(),
		width:         80,
		height:        24,
	}
}

// ── Init ─────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// ── Update ───────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		w := msg.Width - 6
		m.categoryList.SetSize(w, clamp(msg.Height-10, 6, 12))
		m.algorithmList.SetSize(w, clamp(msg.Height-10, 8, 20))
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.state {
		case StateCategorySelect:
			return m.handleCategorySelect(msg)
		case StateAlgorithmSelect:
			return m.handleAlgorithmSelect(msg)
		case StateInputMode:
			return m.handleInputMode(msg)
		case StateTextInput:
			return m.handleTextInput(msg)
		case StateHashing:
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

// ── Key handlers ─────────────────────────────────────────────────────────────

func (m Model) handleCategorySelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "enter":
		sel, ok := m.categoryList.SelectedItem().(categoryItem)
		if !ok {
			return m, nil
		}
		m.selectedCategory = sel.category
		m.algorithmList = makeAlgorithmList(sel.category)
		m.algorithmList.SetSize(m.width-6, clamp(m.height-10, 8, 20))
		m.state = StateAlgorithmSelect
		return m, nil
	default:
		var cmd tea.Cmd
		m.categoryList, cmd = m.categoryList.Update(msg)
		return m, cmd
	}
}

func (m Model) handleAlgorithmSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// While the user is actively typing a filter, forward everything to the list.
	if m.algorithmList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.algorithmList, cmd = m.algorithmList.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc":
		// If a filter is applied, ESC clears it; otherwise go back.
		if m.algorithmList.FilterState() != list.Unfiltered {
			var cmd tea.Cmd
			m.algorithmList, cmd = m.algorithmList.Update(msg)
			return m, cmd
		}
		m.state = StateCategorySelect
		return m, nil
	case "enter":
		sel, ok := m.algorithmList.SelectedItem().(algorithmItem)
		if !ok {
			return m, nil
		}
		m.selectedAlgo = sel.alg
		m.opts.Algorithm = getAlgorithmKey(sel.alg.Name)
		m.state = StateInputMode
		return m, nil
	default:
		var cmd tea.Cmd
		m.algorithmList, cmd = m.algorithmList.Update(msg)
		return m, cmd
	}
}

func (m Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc":
		m.state = StateAlgorithmSelect
	case "s", "1":
		m.inputMode = InputModeString
		m.textInput.Reset()
		m.textInput.Focus()
		m.state = StateTextInput
		return m, textinput.Blink
	case "f", "2":
		m.inputMode = InputModeFile
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
		m.state = StateCategorySelect
		m.results = nil
		m.err = nil
		m.textInput.Reset()
		m.files = nil
		m.categoryList = makeCategoryList()
		m.categoryList.SetSize(m.width-6, clamp(m.height-10, 6, 12))
	case "n":
		m.results = nil
		m.err = nil
		m.textInput.Reset()
		m.state = StateInputMode
	}
	return m, nil
}

// ── Async commands ────────────────────────────────────────────────────────────

func (m Model) doHashString(input string) tea.Cmd {
	return func() tea.Msg {
		r := hasher.HashString(input, m.opts)
		return hashCompleteMsg{results: []hasher.Result{r}}
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

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	var s strings.Builder
	switch m.state {
	case StateCategorySelect:
		s.WriteString(m.viewCategorySelect())
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

func (m Model) viewCategorySelect() string {
	var s strings.Builder
	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(LogoAccent.Render(" ⟡"))
	s.WriteString("\n")
	s.WriteString(SubtitleStyle.Render("compute cryptographic hashes for strings & files"))
	s.WriteString("\n\n")
	s.WriteString(m.categoryList.View())
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/↓  navigate   enter  select   q  quit"))
	return s.String()
}

func (m Model) viewAlgorithmSelect() string {
	var s strings.Builder
	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(LogoAccent.Render(" ⟡ "))
	s.WriteString(LabelStyle.Render(strings.ToUpper(m.selectedCategory.String())))
	s.WriteString("\n\n")
	s.WriteString(m.algorithmList.View())
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/↓  navigate   enter  select   /  filter   esc  back   q  quit"))
	return s.String()
}

func (m Model) viewInputMode() string {
	var s strings.Builder
	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString(LogoAccent.Render(" ⟡ "))
	s.WriteString(LabelStyle.Render(strings.ToUpper(m.selectedAlgo.Name)))
	s.WriteString("\n\n")
	s.WriteString(SubtitleStyle.Render("what do you want to hash?"))
	s.WriteString("\n\n")
	s.WriteString(SelectedStyle.Render("s  hash a string"))
	s.WriteString("\n\n")
	s.WriteString(UnselectedStyle.Render("f  hash a file"))
	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("s  string   f  file   esc  back   q  quit"))
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
	s.WriteString(InputStyle.Render(m.textInput.View()))
	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("enter  hash   esc  back"))
	return s.String()
}

func (m Model) viewHashing() string {
	var s strings.Builder
	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString("\n\n")
	s.WriteString(m.spinner.View())
	s.WriteString("  ")
	s.WriteString(MutedStyle.Render("computing hash..."))
	s.WriteString("\n\n")
	elapsed := time.Since(m.hashStart)
	s.WriteString(DimStyle.Render(fmt.Sprintf("elapsed: %s", elapsed.Round(time.Millisecond))))
	return s.String()
}

func (m Model) viewResults() string {
	var s strings.Builder
	s.WriteString(LogoStyle.Render("hashctl"))
	s.WriteString("  ")
	s.WriteString(SuccessStyle.Render("✓"))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(ErrorStyle.Render("error: " + m.err.Error()))
		s.WriteString("\n")
	} else {
		s.WriteString(LabelStyle.Render(strings.ToUpper(m.selectedAlgo.Name)))
		s.WriteString("\n\n")
		for _, r := range m.results {
			if r.Error != nil {
				s.WriteString(ErrorStyle.Render("✗ "+r.Input))
				s.WriteString("\n")
				s.WriteString(MutedStyle.Render("  "+r.Error.Error()))
				s.WriteString("\n\n")
			} else {
				if r.IsFile {
					s.WriteString(FileStyle.Render("file: "))
					s.WriteString(ValueStyle.Render(r.Input))
				} else {
					s.WriteString(StringStyle.Render("text: "))
					s.WriteString(MutedStyle.Render(`"` + truncate(r.Input, 30) + `"`))
				}
				s.WriteString("\n\n")
				s.WriteString(HashStyle.Render(r.Hash))
				s.WriteString("\n\n")
				s.WriteString(DimStyle.Render(fmt.Sprintf("computed in %s", r.Duration.Round(time.Microsecond))))
				s.WriteString("\n")
			}
		}
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("n  new hash   r  restart   q  quit"))
	return s.String()
}

// ── Helpers ───────────────────────────────────────────────────────────────────

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

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Run starts the TUI application.
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
