package table

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	columnKeySelect = "___select___"
)

var (
	defaultHighlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("#334"))
)

// Model is the main table model.  Create using New().
type Model struct {
	keyMap      KeyMap
	columns     []Column
	headerStyle lipgloss.Style

	rows []Row

	selectableRows bool

	rowCursorIndex int

	focused bool

	highlightStyle lipgloss.Style

	selectedRows []Row

	border Border
}

// New creates a new table ready for further modifications.
func New(columns []Column) Model {
	model := Model{
		columns:        make([]Column, len(columns)),
		highlightStyle: defaultHighlightStyle.Copy(),
		border:         borderDefault,
		keyMap:         DefaultKeyMap(),
	}

	// Do a full deep copy to avoid unexpected edits
	copy(model.columns, columns)

	return model
}

// Init initializes the table per the Bubble Tea architecture.
func (m Model) Init() tea.Cmd {
	return nil
}
