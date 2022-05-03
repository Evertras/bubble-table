package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) moveHighlightUp() {
	m.rowCursorIndex--

	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = len(m.GetVisibleRows()) - 1
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightDown() {
	m.rowCursorIndex++

	if m.rowCursorIndex >= len(m.GetVisibleRows()) {
		m.rowCursorIndex = 0
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) toggleSelect() {
	if !m.selectableRows || len(m.GetVisibleRows()) == 0 {
		return
	}

	rows := make([]Row, len(m.GetVisibleRows()))
	copy(rows, m.GetVisibleRows())

	rows[m.rowCursorIndex].selected = !rows[m.rowCursorIndex].selected

	m.rows = rows

	m.selectedRows = []Row{}

	for _, row := range m.GetVisibleRows() {
		if row.selected {
			m.selectedRows = append(m.selectedRows, row)
		}
	}
}

func (m Model) updateFilterTextInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.FilterBlur) {
			m.filterTextInput.Blur()
		}
	}
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)
	m.pageFirst()

	return m, cmd
}

func (m *Model) handleKeypress(msg tea.KeyMsg) {
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

	case key.Matches(msg, m.keyMap.PageFirst):
		m.pageFirst()

	case key.Matches(msg, m.keyMap.PageLast):
		m.pageLast()

	case key.Matches(msg, m.keyMap.Filter):
		m.filterTextInput.Focus()

	case key.Matches(msg, m.keyMap.FilterClear):
		m.filterTextInput.Reset()
	}
}

// Update responds to input from the user or other messages from Bubble Tea.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	if m.filterTextInput.Focused() {
		return m.updateFilterTextInput(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeypress(msg)
	}

	return m, nil
}
