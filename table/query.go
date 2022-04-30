package table

// GetColumnSorting returns the current sorting rules for the table as a list of
// SortColumns, which are applied from first to last.  This means that data will
// be grouped by the later elements in the list.  The returned list is a copy
// and modifications will have no effect.
func (m *Model) GetColumnSorting() []SortColumn {
	c := make([]SortColumn, len(m.sortOrder))

	copy(c, m.sortOrder)

	return c
}

// GetCanFilter returns true if the table enables filtering at all.  This does
// not say whether a filter is currently active, only that the feature is enabled.
func (m *Model) GetCanFilter() bool {
	return m.filtered
}

// GetIsFilterActive returns true if the table is currently being filtered.  This
// does not say whether the table CAN be filtered, only whether or not a filter
// is actually currently being applied.
func (m *Model) GetIsFilterActive() bool {
	return m.filterTextInput.Value() != ""
}

// GetCurrentFilter returns the current filter text being applied, or an empty
// string if none is applied.
func (m *Model) GetCurrentFilter() string {
	return m.filterTextInput.Value()
}
