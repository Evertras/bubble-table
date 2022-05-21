package table

import (
	"github.com/charmbracelet/bubbles/textinput"
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
	baseStyle      lipgloss.Style
	highlightStyle lipgloss.Style
	headerStyle    lipgloss.Style
	border         Border
	selectedText   string
	unselectedText string

	// Footers
	footerVisible bool
	staticFooter  string

	// Pagination
	pageSize    int
	currentPage int

	// Sorting, where a stable sort is applied from first element to last so
	// that elements are grouped by the later elements.
	sortOrder []SortColumn

	// Filter
	filtered        bool
	filterTextInput textinput.Model

	// For flex columns
	targetTotalWidth int

	// Internal cached calculations for reference
	totalWidth int
}

// New creates a new table ready for further modifications.
func New(columns []Column) Model {
	filterInput := textinput.New()
	filterInput.Prompt = "/"
	model := Model{
		columns:        make([]Column, len(columns)),
		highlightStyle: defaultHighlightStyle.Copy(),
		border:         borderDefault,
		footerVisible:  true,
		keyMap:         DefaultKeyMap(),

		selectedText:   "[x]",
		unselectedText: "[ ]",

		filterTextInput: filterInput,
		baseStyle:       lipgloss.NewStyle().Align(lipgloss.Right),
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

// GetVisibleRows return sorted and filtered rows.
func (m Model) GetVisibleRows() []Row {
	rows := make([]Row, len(m.rows))
	copy(rows, m.rows)
	if m.filtered {
		rows = m.getFilteredRows(rows)
	}
	rows = getSortedRows(m.sortOrder, rows)

	return rows
}
