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
	// Data
	columns []Column
	rows    []Row

	// Interaction
	focused        bool
	keyMap         KeyMap
	selectableRows bool
	selectedRows   []Row
	rowCursorIndex int

	// Styles
	highlightStyle lipgloss.Style
	headerStyle    lipgloss.Style
	border         Border

	// Footers
	staticFooter string

	// Internal cached calculations for reference
	totalWidth int
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

	model.recalculateWidth()

	return model
}

// Init initializes the table per the Bubble Tea architecture.
func (m Model) Init() tea.Cmd {
	return nil
}
