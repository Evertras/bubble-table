package table

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Searchable enables the built-in search keys and footer input.  Unlike
// filtering, search does not hide rows; it moves the highlighted row to the
// next or previous match.
func (m Model) Searchable(searchable bool) Model {
	m.searchable = searchable

	if m.minimumHeight > 0 || m.targetHeight != 0 {
		m.recalculateHeight()
	}

	if m.targetHeight != 0 {
		m.pageStartIndices = nil
	}

	return m
}

// WithSearchFunc sets the function used to decide whether a row matches the
// search query.  If nil, the default is used: case-insensitive "contains"
// matching over all columns.
func (m Model) WithSearchFunc(isMatch FilterFunc) Model {
	if isMatch == nil {
		isMatch = searchFuncContains
	}

	m.searchFunc = isMatch

	return m
}

// WithFuzzySearch enables fuzzy searching for the table.
func (m Model) WithFuzzySearch() Model {
	return m.WithSearchFunc(searchFuncFuzzy)
}

// WithSearchInputValue sets the search query without moving the highlighted
// row.  Call SearchNext or SearchPrev afterwards to jump to a match.
func (m Model) WithSearchInputValue(value string) Model {
	m.searchTextInput.SetValue(value)
	m.searchTextInput.Blur()

	return m
}

// SearchNext highlights the next visible row after the current one that matches
// the search query, wrapping around.  Does nothing if nothing matches.
func (m Model) SearchNext() Model {
	if index := m.findSearchMatch(m.rowCursorIndex+1, 1); index >= 0 {
		return m.WithHighlightedRow(index)
	}

	return m
}

// SearchPrev highlights the previous visible row before the current one that
// matches the search query, wrapping around.  Does nothing if nothing matches.
func (m Model) SearchPrev() Model {
	if index := m.findSearchMatch(m.rowCursorIndex-1, -1); index >= 0 {
		return m.WithHighlightedRow(index)
	}

	return m
}

func (m *Model) handleSearchKeypress(msg tea.KeyPressMsg) {
	if !m.searchable {
		return
	}

	if key.Matches(msg, m.keyMap.Search) {
		m.searchOriginIndex = m.rowCursorIndex
		m.searchTextInput.Focus()
	}

	if key.Matches(msg, m.keyMap.SearchNext) {
		*m = m.SearchNext()
	}

	if key.Matches(msg, m.keyMap.SearchPrev) {
		*m = m.SearchPrev()
	}

	// Shares the filter's clear key so one esc clears both
	if key.Matches(msg, m.keyMap.FilterClear) {
		m.searchTextInput.Reset()
	}
}

func (m Model) isSearchMatch(row Row) bool {
	query := m.searchTextInput.Value()
	if query == "" {
		return false
	}

	return m.searchFunc(FilterFuncInput{
		Columns:        m.columns,
		Row:            row,
		Filter:         query,
		GlobalMetadata: m.metadata,
	})
}

// findSearchMatch returns the index of the first visible row matching the
// search, starting at start and moving by step with wraparound, or -1.
func (m *Model) findSearchMatch(start, step int) int {
	rows := m.GetVisibleRows()
	count := len(rows)

	for i := range count {
		index := ((start+i*step)%count + count) % count

		if m.isSearchMatch(rows[index]) {
			return index
		}
	}

	return -1
}

// updateSearchTextInput searches incrementally: each change jumps to the first
// match at or after the row that was highlighted when search started.
func (m Model) updateSearchTextInput(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok && key.Matches(msg, m.keyMap.SearchBlur) {
		m.searchTextInput.Blur()

		return m, nil
	}

	previousRowIndex := m.rowCursorIndex
	previousQuery := m.searchTextInput.Value()

	var cmd tea.Cmd
	m.searchTextInput, cmd = m.searchTextInput.Update(msg)

	if m.searchTextInput.Value() == previousQuery {
		return m, cmd
	}

	target := m.searchOriginIndex
	if index := m.findSearchMatch(m.searchOriginIndex, 1); index >= 0 {
		target = index
	}

	m = m.WithHighlightedRow(target)

	if m.rowCursorIndex != previousRowIndex {
		m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousRowIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}

	return m, cmd
}
