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
		if m.filterFunc != nil {
			if m.filterFunc(row, filterInputValue) {
				filteredRows = append(filteredRows, row)
			}
		} else {
			if isRowMatched(m.columns, row, filterInputValue) {
				filteredRows = append(filteredRows, row)
			}
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

// NewFuzzyFilter returns a filterFunc that performs case-insensitive fuzzy
// matching (subsequence) over the concatenation of all filterable column values.
// Example wiring:
//
//	m.filterFunc = NewFuzzyFilter(m.columns)
func NewFuzzyFilter(columns []Column) func(Row, string) bool {
	return func(row Row, filter string) bool {
		filter = strings.TrimSpace(filter)
		if filter == "" {
			return true
		}

		// Concatenate all filterable values for this row into one string
		var b strings.Builder
		for _, col := range columns {
			if !col.filterable {
				continue
			}
			if v, ok := row.Data[col.key]; ok {
				// Unwrap StyledCell if present
				switch vv := v.(type) {
				case StyledCell:
					v = vv.Data
				}

				switch vv := v.(type) {
				case string:
					b.WriteString(vv)
				case fmt.Stringer:
					b.WriteString(vv.String())
				default:
					b.WriteString(fmt.Sprintf("%v", v))
				}
				b.WriteByte(' ')
			}
		}

		haystack := strings.ToLower(b.String())
		if haystack == "" {
			return false
		}

		// Support multi-token filters: "acme stl" must fuzzy-match both tokens
		for _, token := range strings.Fields(strings.ToLower(filter)) {
			if !fuzzySubsequenceMatch(haystack, token) {
				return false
			}
		}
		return true
	}
}

// fuzzySubsequenceMatch returns true if all runes in needle appear in order
// within haystack (not necessarily contiguously). Case must be normalized by caller.
func fuzzySubsequenceMatch(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	hi, ni := 0, 0
	hr := []rune(haystack)
	nr := []rune(needle)

	for hi < len(hr) && ni < len(nr) {
		if hr[hi] == nr[ni] {
			ni++
		}
		hi++
	}
	return ni == len(nr)
}
