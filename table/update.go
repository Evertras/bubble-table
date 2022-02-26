package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) moveHighlightUp() {
	m.rowCursorIndex--

	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = len(m.rows) - 1
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightDown() {
	m.rowCursorIndex++

	if m.rowCursorIndex >= len(m.rows) {
		m.rowCursorIndex = 0
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) toggleSelect() {
	if !m.selectableRows || len(m.rows) == 0 {
		return
	}

	rows := make([]Row, len(m.sortedRows))
	copy(rows, m.sortedRows)

	rows[m.rowCursorIndex].selected = !rows[m.rowCursorIndex].selected

	m.sortedRows = rows

	m.selectedRows = []Row{}

	for _, row := range m.sortedRows {
		if row.selected {
			m.selectedRows = append(m.selectedRows, row)
		}
	}
}

// Update responds to input from the user or other messages from Bubble Tea.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.RowDown):
			m.moveHighlightDown()

		case key.Matches(msg, m.keyMap.RowUp):
			m.moveHighlightUp()

		case key.Matches(msg, m.keyMap.RowSelectToggle):
			m.toggleSelect()

		case key.Matches(msg, m.keyMap.PageDown):
			m.pageDown()

		case key.Matches(msg, m.keyMap.PageUp):
			m.pageUp()
		}
	}

	return m, nil
}
