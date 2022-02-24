package table

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strings"
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
	selectedText   string
	unselectedText string

	// Footers
	staticFooter string

	// Pagination
	pageSize    int
	currentPage int

	// Sorting
	sortOrder []sortColumn

	// Filter
	filtered        bool
	filterTextInput textinput.Model

	// Internal cached calculations for reference
	totalWidth int
}

// New creates a new table ready for further modifications.
func New(columns []Column) Model {
	filterInput := textinput.New()
	filterInput.Prompt = ""
	model := Model{
		columns:        make([]Column, len(columns)),
		highlightStyle: defaultHighlightStyle.Copy(),
		border:         borderDefault,
		keyMap:         DefaultKeyMap(),

		selectedText:   "[x]",
		unselectedText: "[ ]",

		filterTextInput: filterInput,
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

func (m Model) isRowMatched(row Row) bool {
	for _, column := range m.columns {
		if column.Filterable {
			data, ok := row.Data[column.Key]
			if !ok {
				continue
			}
			if strings.Contains(strings.ToLower(data.(string)), strings.ToLower(m.filterTextInput.Value())) {
				return true
			}
		}
	}
	return false
}

func (m Model) getFilteredRows(rows []Row) []Row {
	if !m.filtered || m.filterTextInput.Value() == "" {
		return rows
	}

	filteredRows := make([]Row, 0)

	for _, row := range rows {
		if m.isRowMatched(row) {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}

// GetRows return sorted and filtered rows
func (m Model) GetRows() []Row {
	rows := make([]Row, len(m.rows))
	copy(rows, m.rows)
	if m.filtered {
		rows = m.getFilteredRows(rows)
	}
	rows = m.getSortedRows(rows)
	return rows
}

func (m Model) getSortedRows(rows []Row) []Row {
	var sortedRows []Row
	if len(m.sortOrder) == 0 {
		sortedRows = rows

		return sortedRows
	}

	sortedRows = make([]Row, len(rows))
	copy(sortedRows, rows)

	for _, byColumn := range m.sortOrder {
		sorted := &sortableTable{
			rows:     sortedRows,
			byColumn: byColumn,
		}

		sort.Stable(sorted)

		sortedRows = sorted.rows
	}
	return sortedRows
}
