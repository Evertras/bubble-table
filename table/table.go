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

	border Border
}

func New(headers []Header) Model {
	m := Model{
		headers:        make([]Header, len(headers)),
		highlightStyle: defaultHighlightStyle.Copy(),
		border:         borderDefault,
	}

	// Do a full deep copy to avoid unexpected edits
	copy(m.headers, headers)

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

	// TODO: Better way to do this without pointers/nil?  Or should it be nil?
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
	numHeaders := len(m.headers)
	hasRows := len(m.rows) > 0

	// Safety valve for empty tables
	if numHeaders == 0 {
		return ""
	}

	body := strings.Builder{}

	headerStrings := []string{}

	var (
		headerStyleLeft  lipgloss.Style
		headerStyleInner lipgloss.Style
		headerStyleRight lipgloss.Style
	)

	if numHeaders == 1 {
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

	for i, header := range m.headers {
		headerSection := fmt.Sprintf(header.fmtString, header.Title)
		var borderStyle lipgloss.Style

		if i == 0 {
			borderStyle = headerStyleLeft
		} else if i < len(m.headers)-1 {
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
