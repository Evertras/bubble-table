package table

import (
	"fmt"
	"strings"
)

// FilterFuncInput is the input to a FilterFunc. It's a struct so we can add more things later
// without breaking compatibility.
type FilterFuncInput struct {
	// Columns is a list of the columns of the table
	Columns []Column

	// Row is the row that's being considered for filtering
	Row Row

	// GlobalMetadata is an arbitrary set of metadata from the table set by WithGlobalMetadata
	GlobalMetadata map[string]any

	// Filter is the filter string input to consider
	Filter string
}

// FilterFunc takes a FilterFuncInput and returns true if the row should be visible,
// or false if the row should be hidden.
type FilterFunc func(FilterFuncInput) bool

func (m Model) getFilteredRows(rows []Row) []Row {
	filterInputValue := m.filterTextInput.Value()
	if !m.filtered || filterInputValue == "" {
		return rows
	}

	filteredRows := make([]Row, 0)

	for _, row := range rows {
		var availableFilterFunc FilterFunc

		if m.filterFunc != nil {
			availableFilterFunc = m.filterFunc
		} else {
			availableFilterFunc = filterFuncContains
		}

		if availableFilterFunc(FilterFuncInput{
			Columns:        m.columns,
			Row:            row,
			Filter:         filterInputValue,
			GlobalMetadata: m.metadata,
		}) {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}

// filterFuncContains returns a filterFunc that performs case-insensitive
// "contains" matching over all filterable columns in a row.
func filterFuncContains(input FilterFuncInput) bool {
	if input.Filter == "" {
		return true
	}

	checkedAny := false

	filterLower := strings.ToLower(input.Filter)

	for _, column := range input.Columns {
		if !column.filterable {
			continue
		}

		checkedAny = true

		data, ok := input.Row.Data[column.key]

		if !ok {
			continue
		}

		// Extract internal StyledCell data
		switch dataV := data.(type) {
		case StyledCell:
			data = dataV.Data
		}

		var target string
		switch dataV := data.(type) {
		case string:
			target = dataV

		case fmt.Stringer:
			target = dataV.String()

		default:
			target = fmt.Sprintf("%v", data)
		}

		if strings.Contains(strings.ToLower(target), filterLower) {
			return true
		}
	}

	return !checkedAny
}

// filterFuncFuzzy returns a filterFunc that performs case-insensitive fuzzy
// matching (subsequence) over the concatenation of all filterable column values.
func filterFuncFuzzy(input FilterFuncInput) bool {
	filter := strings.TrimSpace(input.Filter)
	if filter == "" {
		return true
	}

	var builder strings.Builder
	for _, col := range input.Columns {
		if !col.filterable {
			continue
		}
		value, ok := input.Row.Data[col.key]
		if !ok {
			continue
		}
		if sc, ok := value.(StyledCell); ok {
			value = sc.Data
		}
		builder.WriteString(fmt.Sprint(value)) // uses Stringer if implemented
		builder.WriteByte(' ')
	}

	haystack := strings.ToLower(builder.String())
	if haystack == "" {
		return false
	}

	for _, token := range strings.Fields(strings.ToLower(filter)) {
		if !fuzzySubsequenceMatch(haystack, token) {
			return false
		}
	}

	return true
}

// fuzzySubsequenceMatch returns true if all runes in needle appear in order
// within haystack (not necessarily contiguously). Case must be normalized by caller.
func fuzzySubsequenceMatch(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	haystackIndex, needleIndex := 0, 0
	haystackRunes := []rune(haystack)
	needleRunes := []rune(needle)

	for haystackIndex < len(haystackRunes) && needleIndex < len(needleRunes) {
		if haystackRunes[haystackIndex] == needleRunes[needleIndex] {
			needleIndex++
		}
		haystackIndex++
	}

	return needleIndex == len(needleRunes)
}
