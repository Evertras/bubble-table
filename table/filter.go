package table

import (
	"fmt"
	"strings"
)

func (m Model) getFilteredRows(rows []Row) []Row {
	filterInputValue := m.filterTextInput.Value()
	if !m.filtered || filterInputValue == "" {
		return rows
	}

	filteredRows := make([]Row, 0)

	for _, row := range rows {
		var availableFilterFunc func([]Column, Row, string) bool

		if m.filterFunc != nil {
			availableFilterFunc = m.filterFunc
		} else {
			availableFilterFunc = isRowMatched
		}

		if availableFilterFunc(m.columns, row, filterInputValue) {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}

func isRowMatched(columns []Column, row Row, filter string) bool {
	if filter == "" {
		return true
	}

	checkedAny := false

	filterLower := strings.ToLower(filter)

	for _, column := range columns {
		if !column.filterable {
			continue
		}

		checkedAny = true

		data, ok := row.Data[column.key]

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

// newFuzzyFilter returns a filterFunc that performs case-insensitive fuzzy
// matching (subsequence) over the concatenation of all filterable column values.
func newFuzzyFilter(columns []Column, row Row, filter string) bool {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return true
	}

	var builder strings.Builder
	for _, col := range columns {
		if !col.filterable {
			continue
		}
		value, ok := row.Data[col.key]
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
