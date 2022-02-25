package table

import "strings"

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

	for _, column := range columns {
		if column.Filterable {
			data, ok := row.Data[column.Key]
			if !ok {
				continue
			}
			if strings.Contains(strings.ToLower(data.(string)), strings.ToLower(filter)) {
				return true
			}
		}
	}
	return false
}
