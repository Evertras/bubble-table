package table

import (
	"fmt"
	"sort"
)

type sortDirection int

const (
	sortDirectionAsc sortDirection = iota
	sortDirectionDesc
)

type sortColumn struct {
	columnKey string
	direction sortDirection
}

func (m Model) SortByAsc(columnKey string) Model {
	m.sortOrder = []sortColumn{
		{
			columnKey: columnKey,
			direction: sortDirectionAsc,
		},
	}

	m.updateSortedRows()

	return m
}

func (m Model) SortByDesc(columnKey string) Model {
	m.sortOrder = []sortColumn{
		{
			columnKey: columnKey,
			direction: sortDirectionDesc,
		},
	}

	m.updateSortedRows()

	return m
}

func (m Model) ThenSortByAsc(columnKey string) Model {
	m.sortOrder = append([]sortColumn{
		{
			columnKey: columnKey,
			direction: sortDirectionAsc,
		},
	}, m.sortOrder...)

	m.updateSortedRows()

	return m
}

func (m Model) ThenSortByDesc(columnKey string) Model {
	m.sortOrder = append([]sortColumn{
		{
			columnKey: columnKey,
			direction: sortDirectionDesc,
		},
	}, m.sortOrder...)

	m.updateSortedRows()

	return m
}

type sortableTable struct {
	rows     []Row
	byColumn sortColumn
}

func (s *sortableTable) Len() int {
	return len(s.rows)
}

func (s *sortableTable) Swap(i, j int) {
	old := s.rows[i]
	s.rows[i] = s.rows[j]
	s.rows[j] = old
}

func (s *sortableTable) extractString(i int, column string) string {
	iData, exists := s.rows[i].Data[column]

	if !exists {
		return ""
	}

	switch iData := iData.(type) {
	case StyledCell:
		return fmt.Sprintf("%v", iData.Data)

	case string:
		return iData

	default:
		return fmt.Sprintf("%v", iData)
	}
}

func (s *sortableTable) extractNumber(i int, column string) (float64, bool) {
	iData, exists := s.rows[i].Data[column]

	if !exists {
		return 0, false
	}

	return asNumber(iData)
}

func (s *sortableTable) Less(first, second int) bool {
	firstNum, firstNumIsValid := s.extractNumber(first, s.byColumn.columnKey)
	secondNum, secondNumIsValid := s.extractNumber(second, s.byColumn.columnKey)

	if firstNumIsValid && secondNumIsValid {
		if s.byColumn.direction == sortDirectionAsc {
			return firstNum < secondNum
		}

		return firstNum > secondNum
	}

	firstVal := s.extractString(first, s.byColumn.columnKey)
	secondVal := s.extractString(second, s.byColumn.columnKey)

	if s.byColumn.direction == sortDirectionAsc {
		return firstVal < secondVal
	}

	return firstVal > secondVal
}

func (m *Model) updateSortedRows() {
	if len(m.sortOrder) == 0 {
		m.sortedRows = m.rows

		return
	}

	m.sortedRows = make([]Row, len(m.rows))
	copy(m.sortedRows, m.rows)

	for _, byColumn := range m.sortOrder {
		sorted := &sortableTable{
			rows:     m.sortedRows,
			byColumn: byColumn,
		}

		sort.Stable(sorted)

		m.sortedRows = sorted.rows
	}
}
