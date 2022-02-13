package table

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MovementMode int

const (
	columnKeySelect = "___select___"
)

var (
	defaultHighlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("#334"))
)

type Model struct {
	headers     []Header
	headerStyle lipgloss.Style

	rows []Row

	selectableRows bool

	rowCursorIndex int

	focused bool

	highlightStyle lipgloss.Style

	selectedRows []Row
}

func New(headers []Header) Model {
	m := Model{
		headers:        make([]Header, len(headers)),
		highlightStyle: defaultHighlightStyle.Copy(),
	}

	// Do a full deep copy to avoid unexpected edits
	copy(m.headers, headers)
	for i, header := range m.headers {
		m.headers[i].Style = header.Style.Copy()
	}

	return m
}

func (m Model) HeaderStyle(style lipgloss.Style) Model {
	m.headerStyle = style.Copy()
	return m
}

func (m Model) WithRows(rows []Row) Model {
	m.rows = rows
	return m
}

func (m Model) SelectableRows(selectable bool) Model {
	m.selectableRows = selectable

	hasSelectHeader := m.headers[0].Key == columnKeySelect

	if hasSelectHeader != selectable {
		if selectable {
			m.headers = append([]Header{
				NewHeader(columnKeySelect, "[x]", 3),
			}, m.headers...)
		} else {
			m.headers = m.headers[1:]
		}
	}

	return m
}

func (m Model) HighlightedRow() Row {
	if len(m.rows) > 0 {
		return m.rows[m.rowCursorIndex]
	}

	// This shouldn't really happen... better indication?
	return Row{}
}

func (m Model) SelectedRows() []Row {
	return m.selectedRows
}

func (m Model) HighlightStyle(style lipgloss.Style) Model {
	m.highlightStyle = style
	return m
}

func (m Model) Focused(focused bool) Model {
	m.focused = focused
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

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

func (m Model) View() string {
	body := strings.Builder{}

	headerStrings := []string{}

	for i, header := range m.headers {
		headerSection := fmt.Sprintf(header.fmtString, header.Title)
		borderStyle := m.headerStyle.Copy()

		if i == 0 {
			borderStyle = borderStyle.BorderStyle(borderHeaderFirst)
		} else if i < len(m.headers)-1 {
			borderStyle = borderStyle.BorderStyle(borderHeaderMiddle).BorderTop(true).BorderBottom(true).BorderRight(true)
		} else {
			borderStyle = borderStyle.BorderStyle(borderHeaderLast).BorderTop(true).BorderBottom(true).BorderRight(true)
		}

		headerStrings = append(headerStrings, borderStyle.Render(header.Style.Render(headerSection)))
	}

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...))

	body.WriteString("\n")

	for i := range m.rows {
		body.WriteString(m.renderRow(i))
		body.WriteString("\n")
	}

	return body.String()
}
