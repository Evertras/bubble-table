package table

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	columnKeySelect = "___select___"
)

var (
	defaultHighlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("#334"))
)

// Model is the main table model.  Create using New()
type Model struct {
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

// New creates a new table ready for further modifications
func New(columns []Column) Model {
	m := Model{
		columns:        make([]Column, len(columns)),
		highlightStyle: defaultHighlightStyle.Copy(),
		border:         borderDefault,
	}

	// Do a full deep copy to avoid unexpected edits
	copy(m.columns, columns)

	return m
}

// HeaderStyle sets the style to apply to the header text, such as color or bold
func (m Model) HeaderStyle(style lipgloss.Style) Model {
	m.headerStyle = style.Copy()
	return m
}

// WithRows sets the rows to show as data in the table
func (m Model) WithRows(rows []Row) Model {
	m.rows = rows
	return m
}

// SelectableRows sets whether or not rows are selectable.  If set, adds a column
// in the front that acts as a checkbox and responds to controls if Focused
func (m Model) SelectableRows(selectable bool) Model {
	m.selectableRows = selectable

	hasSelectColumn := m.columns[0].Key == columnKeySelect

	if hasSelectColumn != selectable {
		if selectable {
			m.columns = append([]Column{
				NewColumn(columnKeySelect, "[x]", 3),
			}, m.columns...)
		} else {
			m.columns = m.columns[1:]
		}
	}

	return m
}

// HighlightedRow returns the full Row that's currently highlighted by the user
func (m Model) HighlightedRow() Row {
	if len(m.rows) > 0 {
		return m.rows[m.rowCursorIndex]
	}

	// TODO: Better way to do this without pointers/nil?  Or should it be nil?
	return Row{}
}

// SelectedRows returns all rows that have been set as selected by the user
func (m Model) SelectedRows() []Row {
	return m.selectedRows
}

// HighlightStyle sets a custom style to use when the row is being highlighted
// by the cursor
func (m Model) HighlightStyle(style lipgloss.Style) Model {
	m.highlightStyle = style
	return m
}

// Focused allows the table to show highlighted rows and take in controls of
// up/down/space/etc to let the user navigate the table and interact with it
func (m Model) Focused(focused bool) Model {
	m.focused = focused
	return m
}

// Init initializes the table per the Bubble Tea architecture
func (m Model) Init() tea.Cmd {
	return nil
}

// Update responds to input from the user or other messages from Bubble Tea
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			m.rowCursorIndex++

			if m.rowCursorIndex >= len(m.rows) {
				m.rowCursorIndex = 0
			}

		case "up", "k":
			m.rowCursorIndex--

			if m.rowCursorIndex < 0 {
				m.rowCursorIndex = len(m.rows) - 1
			}

		case " ", "enter":
			if !m.selectableRows || len(m.rows) == 0 {
				break
			}

			rows := make([]Row, len(m.rows))
			copy(rows, m.rows)

			rows[m.rowCursorIndex].selected = !rows[m.rowCursorIndex].selected

			m.rows = rows

			m.selectedRows = []Row{}

			for _, row := range m.rows {
				if row.selected {
					m.selectedRows = append(m.selectedRows, row)
				}
			}
		}

	}

	return m, nil
}

// View renders the table.  It does not end in a newline, so that it can be
// composed with other elements more consistently.
func (m Model) View() string {
	numColumns := len(m.columns)
	hasRows := len(m.rows) > 0

	// Safety valve for empty tables
	if numColumns == 0 {
		return ""
	}

	body := strings.Builder{}

	headerStrings := []string{}

	var (
		headerStyleLeft  lipgloss.Style
		headerStyleInner lipgloss.Style
		headerStyleRight lipgloss.Style
	)

	if numColumns == 1 {
		if hasRows {
			headerStyleLeft = m.border.styleSingleColumnTop
		} else {
			headerStyleLeft = m.border.styleSingleCell
		}

		headerStyleLeft = headerStyleLeft.Copy().Inherit(m.headerStyle)
	} else {
		if hasRows {
			headerStyleLeft = m.border.styleMultiTopLeft
			headerStyleInner = m.border.styleMultiTop
			headerStyleRight = m.border.styleMultiTopRight
		} else {
			headerStyleLeft = m.border.styleSingleRowLeft
			headerStyleInner = m.border.styleSingleRowInner
			headerStyleRight = m.border.styleSingleRowRight
		}

		headerStyleLeft = headerStyleLeft.Copy().Inherit(m.headerStyle)
		headerStyleInner = headerStyleInner.Copy().Inherit(m.headerStyle)
		headerStyleRight = headerStyleRight.Copy().Inherit(m.headerStyle)
	}

	for i, header := range m.columns {
		headerSection := fmt.Sprintf(header.fmtString, header.Title)
		var borderStyle lipgloss.Style

		if i == 0 {
			borderStyle = headerStyleLeft
		} else if i < len(m.columns)-1 {
			borderStyle = headerStyleInner
		} else {
			borderStyle = headerStyleRight
		}

		headerStrings = append(headerStrings, borderStyle.Render(headerSection))
	}

	headerBlock := lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...)

	rowStrs := []string{headerBlock}
	for i := range m.rows {
		rowStrs = append(rowStrs, m.renderRow(i))
	}

	body.WriteString(lipgloss.JoinVertical(lipgloss.Left, rowStrs...))

	return body.String()
}
