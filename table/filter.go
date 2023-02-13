package table

import (
	"fmt"
	"strings"
)

func (m Model) getFilteredRows(rows []Row) []Row {
	if !m.filtered || m.filterTextInput.Value() == "" {
		return rows
	}

	filteredRows := make([]Row, 0)

	for _, row := range rows {
		if isRowMatched(m.columns, row, m.filterTextInput.Value()) {
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

		target := ""
		switch dataV := data.(type) {
		case string:
			target = dataV

		case fmt.Stringer:
			target = dataV.String()
		}

		if strings.Contains(strings.ToLower(target), strings.ToLower(filter)) {
			return true
		}
	}

	return !checkedAny
}
