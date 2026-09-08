package table

import (
	"cmp"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SortDirection indicates whether a column should sort by ascending or descending.
type SortDirection int

const (
	// SortDirectionAsc indicates the column should be in ascending order.
	SortDirectionAsc SortDirection = iota

	// SortDirectionDesc indicates the column should be in descending order.
	SortDirectionDesc
)

// SortColumn describes which column should be sorted and how.
type SortColumn struct {
	ColumnKey string
	Direction SortDirection
}

// SortByAsc sets the main sorting column to the given key, in ascending order.
// If a previous sort was used, it is replaced by the given column each time
// this function is called.  With the column's default SortType (SortTypeAuto),
// the whole column is compared as numbers if every value is a number, as
// times if every value is a time, or as strings otherwise. See WithSortType
// and WithSortFunc on Column to choose a different comparison per column.
func (m Model) SortByAsc(columnKey string) Model {
	m.sortOrder = []SortColumn{
		{
			ColumnKey: columnKey,
			Direction: SortDirectionAsc,
		},
	}

	m.visibleRowCacheUpdated = false

	return m
}

// SortByDesc sets the main sorting column to the given key, in descending order.
// If a previous sort was used, it is replaced by the given column each time
// this function is called.  With the column's default SortType (SortTypeAuto),
// the whole column is compared as numbers if every value is a number, as
// times if every value is a time, or as strings otherwise. See WithSortType
// and WithSortFunc on Column to choose a different comparison per column.
func (m Model) SortByDesc(columnKey string) Model {
	m.sortOrder = []SortColumn{
		{
			ColumnKey: columnKey,
			Direction: SortDirectionDesc,
		},
	}

	m.visibleRowCacheUpdated = false

	return m
}

// ThenSortByAsc provides a secondary sort after the first, in ascending order.
// Can be chained multiple times, applying to smaller subgroups each time.
func (m Model) ThenSortByAsc(columnKey string) Model {
	m.sortOrder = append([]SortColumn{
		{
			ColumnKey: columnKey,
			Direction: SortDirectionAsc,
		},
	}, m.sortOrder...)

	m.visibleRowCacheUpdated = false

	return m
}

// ThenSortByDesc provides a secondary sort after the first, in descending order.
// Can be chained multiple times, applying to smaller subgroups each time.
func (m Model) ThenSortByDesc(columnKey string) Model {
	m.sortOrder = append([]SortColumn{
		{
			ColumnKey: columnKey,
			Direction: SortDirectionDesc,
		},
	}, m.sortOrder...)

	m.visibleRowCacheUpdated = false

	return m
}

type sortableTable struct {
	rows     []Row
	columns  []Column
	byColumn SortColumn

	autoKindComputed bool
	autoKindValue    autoKind
}

// autoKind is the single comparison metric SortTypeAuto applies across an
// entire column, determined once from all of the column's values rather than
// re-decided per pair. Deciding per pair (the pre-fix behavior) isn't
// transitive: e.g. with 9, 10 (numbers) and "10+" (string) in one column,
// 9 < 10 numerically, 10 < "10+" as strings, but 9 > "10+" as strings, a
// 3-cycle that makes sort.Stable's result depend on input order (see
// https://github.com/Evertras/bubble-table/issues/232).
type autoKind int

const (
	autoKindNumeric autoKind = iota
	autoKindTime
	autoKindString
)

// resolveAutoKind returns the metric SortTypeAuto should use for this
// column, computing and caching it on first use. A column is numeric (or
// time) only if every value present in it is numeric (or time); a column
// with no values, or with any value that doesn't fit one shared kind, falls
// back to string comparison for all of its values.
func (s *sortableTable) resolveAutoKind() autoKind {
	if !s.autoKindComputed {
		s.autoKindValue = detectAutoKind(s.rows, s.byColumn.ColumnKey)
		s.autoKindComputed = true
	}

	return s.autoKindValue
}

func detectAutoKind(rows []Row, columnKey string) autoKind {
	sawValue, allNumeric, allTime := false, true, true

	for i := range rows {
		data, exists := rows[i].Data[columnKey]
		if !exists {
			continue
		}

		sawValue = true

		if _, ok := asNumber(data); !ok {
			allNumeric = false
		}

		if _, ok := asTime(data); !ok {
			allTime = false
		}
	}

	switch {
	case sawValue && allNumeric:
		return autoKindNumeric

	case sawValue && allTime:
		return autoKindTime

	default:
		return autoKindString
	}
}

func (s *sortableTable) column() Column {
	for _, column := range s.columns {
		if column.Key() == s.byColumn.ColumnKey {
			return column
		}
	}

	return Column{}
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

func (s *sortableTable) extractValue(i int, column string) any {
	iData, exists := s.rows[i].Data[column]

	if !exists {
		return nil
	}

	if styled, ok := iData.(StyledCell); ok {
		return styled.Data
	}

	return iData
}

func (s *sortableTable) extractNumber(i int, column string) (float64, bool) {
	iData, exists := s.rows[i].Data[column]

	if !exists {
		return 0, false
	}

	return asNumber(iData)
}

func (s *sortableTable) extractTime(i int, column string) (time.Time, bool) {
	iData, exists := s.rows[i].Data[column]

	if !exists {
		return time.Time{}, false
	}

	return asTime(iData)
}

func (s *sortableTable) Less(first, second int) bool {
	column := s.column()

	if sortFunc := column.SortFunc(); sortFunc != nil {
		return s.lessCustom(first, second, sortFunc)
	}

	if column.SortType() == SortTypeNatural {
		return s.lessNatural(first, second)
	}

	return s.lessAuto(first, second)
}

// directional applies less in ascending order, or with its arguments
// swapped in descending order, so every kind of comparison shares one place
// that knows how SortDirection flips a comparator.
func directional[T any](direction SortDirection, first, second T, less func(a, b T) bool) bool {
	if direction == SortDirectionAsc {
		return less(first, second)
	}

	return less(second, first)
}

func (s *sortableTable) lessCustom(first, second int, sortFunc func(a, b any) bool) bool {
	firstVal := s.extractValue(first, s.byColumn.ColumnKey)
	secondVal := s.extractValue(second, s.byColumn.ColumnKey)

	return directional(s.byColumn.Direction, firstVal, secondVal, sortFunc)
}

func (s *sortableTable) lessNatural(first, second int) bool {
	firstVal := s.extractString(first, s.byColumn.ColumnKey)
	secondVal := s.extractString(second, s.byColumn.ColumnKey)

	return directional(s.byColumn.Direction, firstVal, secondVal, naturalLess)
}

func (s *sortableTable) lessAuto(first, second int) bool {
	switch s.resolveAutoKind() {
	case autoKindNumeric:
		return s.lessNumericKind(first, second)

	case autoKindTime:
		return s.lessTimeKind(first, second)

	case autoKindString:
		return s.lessStringKind(first, second)
	}

	return s.lessStringKind(first, second)
}

func (s *sortableTable) lessNumericKind(first, second int) bool {
	firstNum, firstOk := s.extractNumber(first, s.byColumn.ColumnKey)
	secondNum, secondOk := s.extractNumber(second, s.byColumn.ColumnKey)

	if !firstOk || !secondOk {
		return s.lessWithMissing(firstOk, secondOk)
	}

	return directional(s.byColumn.Direction, firstNum, secondNum, func(a, b float64) bool { return a < b })
}

func (s *sortableTable) lessTimeKind(first, second int) bool {
	firstTime, firstOk := s.extractTime(first, s.byColumn.ColumnKey)
	secondTime, secondOk := s.extractTime(second, s.byColumn.ColumnKey)

	if !firstOk || !secondOk {
		return s.lessWithMissing(firstOk, secondOk)
	}

	return directional(s.byColumn.Direction, firstTime, secondTime, time.Time.Before)
}

func (s *sortableTable) lessStringKind(first, second int) bool {
	firstVal := s.extractString(first, s.byColumn.ColumnKey)
	secondVal := s.extractString(second, s.byColumn.ColumnKey)

	return directional(s.byColumn.Direction, firstVal, secondVal, func(a, b string) bool { return a < b })
}

// lessWithMissing orders a value missing from a row before any present value
// in ascending order (and after, in descending), matching how a missing
// value's "" string representation already sorted relative to real values
// under the old per-pair chain. Both missing counts as equal, preserving
// sort.Stable's original relative order.
func (s *sortableTable) lessWithMissing(firstOk, secondOk bool) bool {
	switch {
	case !firstOk && !secondOk:
		return false

	case !firstOk:
		return s.byColumn.Direction == SortDirectionAsc

	default:
		return s.byColumn.Direction == SortDirectionDesc
	}
}

func getSortedRows(sortOrder []SortColumn, rows []Row, columns []Column) []Row {
	var sortedRows []Row
	if len(sortOrder) == 0 {
		sortedRows = rows

		return sortedRows
	}

	sortedRows = make([]Row, len(rows))
	copy(sortedRows, rows)

	for _, byColumn := range sortOrder {
		sorted := &sortableTable{
			rows:     sortedRows,
			columns:  columns,
			byColumn: byColumn,
		}

		sort.Stable(sorted)

		sortedRows = sorted.rows
	}

	return sortedRows
}

// naturalLess reports whether a should sort before b in natural order: runs
// of ASCII digits are compared numerically rather than character-by-character
// (so "item2" sorts before "item10"), everything else is compared byte-by-byte.
// When two digit runs are numerically equal, the one with fewer leading zeros
// sorts first (so "2" sorts before "02").
func naturalLess(a, b string) bool {
	return naturalCompare(a, b) < 0
}

func naturalCompare(first, second string) int {
	firstIdx, secondIdx := 0, 0

	for firstIdx < len(first) && secondIdx < len(second) {
		if isASCIIDigit(first[firstIdx]) && isASCIIDigit(second[secondIdx]) {
			var runA, runB string

			runA, firstIdx = consumeDigitRun(first, firstIdx)
			runB, secondIdx = consumeDigitRun(second, secondIdx)

			if result := compareDigitRuns(runA, runB); result != 0 {
				return result
			}

			continue
		}

		if first[firstIdx] != second[secondIdx] {
			return cmp.Compare(first[firstIdx], second[secondIdx])
		}

		firstIdx++
		secondIdx++
	}

	return cmp.Compare(len(first)-firstIdx, len(second)-secondIdx)
}

func consumeDigitRun(value string, start int) (string, int) {
	end := start

	for end < len(value) && isASCIIDigit(value[end]) {
		end++
	}

	return value[start:end], end
}

func compareDigitRuns(runA, runB string) int {
	trimmedA := strings.TrimLeft(runA, "0")
	trimmedB := strings.TrimLeft(runB, "0")

	if len(trimmedA) != len(trimmedB) {
		return cmp.Compare(len(trimmedA), len(trimmedB))
	}

	if trimmedA != trimmedB {
		if trimmedA < trimmedB {
			return -1
		}

		return 1
	}

	// Numerically equal: fewer leading zeros sorts first.
	return cmp.Compare(len(runA), len(runB))
}

func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
