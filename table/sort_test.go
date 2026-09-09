package table

import (
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestSortSingleColumnAscAndDesc(t *testing.T) {
	const idColKey = "id"

	// Check mixing types
	type someType string

	rows := []Row{
		NewRow(RowData{idColKey: someType("b")}),
		NewRow(RowData{idColKey: NewStyledCell("c", lipgloss.NewStyle().Bold(true))}),
		NewRow(RowData{idColKey: "a"}),
		// Missing data
		NewRow(RowData{}),
	}

	model := New([]Column{
		NewColumn(idColKey, "ID", 3),
	}).WithRows(rows).SortByAsc(idColKey)

	assertOrder := func(expectedList []string) {
		for index, expected := range expectedList {
			idVal, ok := model.GetVisibleRows()[index].Data[idColKey]

			if expected != "" {
				assert.True(t, ok)
			} else {
				assert.False(t, ok)

				continue
			}

			switch idVal := idVal.(type) {
			case string:
				assert.Equal(t, expected, idVal)

			case someType:
				assert.Equal(t, expected, string(idVal))

			case StyledCell:
				assert.Equal(t, expected, idVal.Data)

			default:
				assert.Fail(t, "Unknown type")
			}
		}
	}

	assert.Len(t, model.GetVisibleRows(), len(rows))
	assertOrder([]string{"", "a", "b", "c"})

	model = model.SortByDesc(idColKey)

	assertOrder([]string{"c", "b", "a", ""})
}

func TestSortSingleColumnIntsAsc(t *testing.T) {
	const idColKey = "id"

	rows := []Row{
		NewRow(RowData{idColKey: 13}),
		NewRow(RowData{idColKey: NewStyledCell(1, lipgloss.NewStyle().Bold(true))}),
		NewRow(RowData{idColKey: 2}),
	}

	model := New([]Column{
		NewColumn(idColKey, "ID", 3),
	}).WithRows(rows).SortByAsc(idColKey)

	assertOrder := func(expectedList []int) {
		for index, expected := range expectedList {
			idVal, ok := model.GetVisibleRows()[index].Data[idColKey]

			assert.True(t, ok)

			switch idVal := idVal.(type) {
			case int:
				assert.Equal(t, expected, idVal)

			case StyledCell:
				assert.Equal(t, expected, idVal.Data)

			default:
				assert.Fail(t, "Unknown type")
			}
		}
	}

	assert.Len(t, model.GetVisibleRows(), len(rows))
	assertOrder([]int{1, 2, 13})
}

func TestSortTwoColumnsAscDescMix(t *testing.T) {
	const (
		nameKey  = "name"
		scoreKey = "score"
	)

	makeRow := func(name string, score int) Row {
		return NewRow(RowData{
			nameKey:  name,
			scoreKey: score,
		})
	}

	model := New([]Column{
		NewColumn(nameKey, "Name", 8),
		NewColumn(scoreKey, "Score", 8),
	}).WithRows([]Row{
		makeRow("c", 50),
		makeRow("a", 75),
		makeRow("b", 101),
		makeRow("a", 100),
	}).SortByAsc(nameKey).ThenSortByDesc(scoreKey)

	assertVals := func(index int, name string, score int) {
		actualName, ok := model.GetVisibleRows()[index].Data[nameKey].(string)
		assert.True(t, ok)

		actualScore, ok := model.GetVisibleRows()[index].Data[scoreKey].(int)
		assert.True(t, ok)

		assert.Equal(t, name, actualName)
		assert.Equal(t, score, actualScore)
	}

	assert.Len(t, model.GetVisibleRows(), 4)

	assertVals(0, "a", 100)
	assertVals(1, "a", 75)
	assertVals(2, "b", 101)
	assertVals(3, "c", 50)

	model = model.SortByDesc(nameKey).ThenSortByAsc(scoreKey)

	assertVals(0, "c", 50)
	assertVals(1, "b", 101)
	assertVals(2, "a", 75)
	assertVals(3, "a", 100)
}

func TestGetSortedRows(t *testing.T) {
	sortColumns := []SortColumn{
		{
			ColumnKey: "cb",
			Direction: SortDirectionDesc,
		},
		{
			ColumnKey: "ca",
			Direction: SortDirectionAsc,
		},
	}
	rows := getSortedRows(sortColumns, []Row{
		NewRow(RowData{
			"ca": "2",
			"cb": "t-1",
		}),
		NewRow(RowData{
			"ca": "1",
			"cb": "t-2",
		}),
		NewRow(RowData{
			"ca": "3",
			"cb": "t-3",
		}),
		NewRow(RowData{
			"ca": "3",
			"cb": "t-2",
		}),
	}, nil)
	assert.Len(t, rows, 4)
	assert.Equal(t, "1", rows[0].Data["ca"])
	assert.Equal(t, "2", rows[1].Data["ca"])
	assert.Equal(t, "3", rows[2].Data["ca"])
	assert.Equal(t, "3", rows[3].Data["ca"])

	assert.Equal(t, "t-2", rows[0].Data["cb"])
	assert.Equal(t, "t-1", rows[1].Data["cb"])
	assert.Equal(t, "t-3", rows[2].Data["cb"])
	assert.Equal(t, "t-2", rows[3].Data["cb"])
}

func TestSortTimeColumnAscDesc(t *testing.T) {
	const timeKey = "time"

	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	makeRow := func(t time.Time) Row { return NewRow(RowData{timeKey: t}) }

	rows := []Row{
		makeRow(base.Add(2 * time.Hour)),
		makeRow(base),
		makeRow(base.Add(time.Hour)),
	}

	model := New([]Column{NewColumn(timeKey, "Time", 20)}).
		WithRows(rows).
		SortByAsc(timeKey)

	getTime := func(i int) time.Time {
		v, ok := model.GetVisibleRows()[i].Data[timeKey].(time.Time)
		assert.True(t, ok)

		return v
	}

	assert.True(t, getTime(0).Equal(base))
	assert.True(t, getTime(1).Equal(base.Add(time.Hour)))
	assert.True(t, getTime(2).Equal(base.Add(2*time.Hour)))

	model = model.SortByDesc(timeKey)

	assert.True(t, getTime(0).Equal(base.Add(2*time.Hour)))
	assert.True(t, getTime(1).Equal(base.Add(time.Hour)))
	assert.True(t, getTime(2).Equal(base))
}

func unwrapStyled(v any) any {
	if styled, ok := v.(StyledCell); ok {
		return styled.Data
	}

	return v
}

func cellAt(model Model, columnKey string, i int) any {
	return unwrapStyled(model.GetVisibleRows()[i].Data[columnKey])
}

func TestSortNaturalType(t *testing.T) {
	const idColKey = "id"

	rows := []Row{
		NewRow(RowData{idColKey: "item10"}),
		NewRow(RowData{idColKey: "item2"}),
		NewRow(RowData{idColKey: "item1"}),
		NewRow(RowData{idColKey: NewStyledCell("item20", lipgloss.NewStyle().Bold(true))}),
	}

	model := New([]Column{
		NewColumn(idColKey, "ID", 8).WithSortType(SortTypeNatural),
	}).WithRows(rows).SortByAsc(idColKey)

	getVal := func(i int) string {
		str, ok := cellAt(model, idColKey, i).(string)
		assert.True(t, ok)

		return str
	}

	assert.Equal(t, "item1", getVal(0))
	assert.Equal(t, "item2", getVal(1))
	assert.Equal(t, "item10", getVal(2))
	assert.Equal(t, "item20", getVal(3))

	model = model.SortByDesc(idColKey)

	assert.Equal(t, "item20", getVal(0))
	assert.Equal(t, "item10", getVal(1))
	assert.Equal(t, "item2", getVal(2))
	assert.Equal(t, "item1", getVal(3))
}

func TestNaturalCompareLeadingZeros(t *testing.T) {
	assert.True(t, naturalLess("2", "02"))
	assert.False(t, naturalLess("02", "2"))
	assert.Equal(t, 0, naturalCompare("02", "02"))
	// Non-digit characters differing at the same position fall back to a
	// plain byte comparison rather than digit-run comparison.
	assert.True(t, naturalLess("apple2", "banana1"))
}

func TestSortCustomFunc(t *testing.T) {
	const idColKey = "id"

	// Sort by string length instead of lexicographically.
	byLength := func(a, b any) bool {
		aStr, ok := a.(string)
		assert.True(t, ok)

		bStr, ok := b.(string)
		assert.True(t, ok)

		return len(aStr) < len(bStr)
	}

	rows := []Row{
		NewRow(RowData{idColKey: "ccc"}),
		NewRow(RowData{idColKey: "a"}),
		NewRow(RowData{idColKey: NewStyledCell("bb", lipgloss.NewStyle().Bold(true))}),
	}

	model := New([]Column{
		NewColumn(idColKey, "ID", 8).WithSortFunc(byLength),
	}).WithRows(rows).SortByAsc(idColKey)

	getVal := func(i int) any {
		return cellAt(model, idColKey, i)
	}

	assert.Equal(t, "a", getVal(0))
	assert.Equal(t, "bb", getVal(1))
	assert.Equal(t, "ccc", getVal(2))

	model = model.SortByDesc(idColKey)

	assert.Equal(t, "ccc", getVal(0))
	assert.Equal(t, "bb", getVal(1))
	assert.Equal(t, "a", getVal(2))
}

func TestSortCustomFuncWithMissingValue(t *testing.T) {
	const idColKey = "id"

	// nil (a missing value) always sorts first.
	byPresence := func(a, _ any) bool {
		return a == nil
	}

	rows := []Row{
		NewRow(RowData{idColKey: "x"}),
		NewRow(RowData{}),
	}

	model := New([]Column{
		NewColumn(idColKey, "ID", 8).WithSortFunc(byPresence),
	}).WithRows(rows).SortByAsc(idColKey)

	_, ok := model.GetVisibleRows()[0].Data[idColKey]
	assert.False(t, ok)
	assert.Equal(t, "x", model.GetVisibleRows()[1].Data[idColKey])
}

func TestSortTypeCustomReportedWhenSortFuncSet(t *testing.T) {
	col := NewColumn("id", "ID", 8)
	assert.Equal(t, SortTypeAuto, col.SortType())

	col = col.WithSortType(SortTypeNatural)
	assert.Equal(t, SortTypeNatural, col.SortType())

	col = col.WithSortFunc(func(_, _ any) bool { return false })
	assert.Equal(t, SortTypeCustom, col.SortType())
}

func TestWithSortTypeCustomDirectlyHasNoEffect(t *testing.T) {
	col := NewColumn("id", "ID", 8).WithSortType(SortTypeCustom)
	assert.Equal(t, SortTypeAuto, col.SortType(), "SortTypeCustom with no comparator shouldn't be reported")

	col = NewColumn("id", "ID", 8).WithSortType(SortTypeNatural).WithSortType(SortTypeCustom)
	assert.Equal(t, SortTypeNatural, col.SortType(), "SortTypeCustom shouldn't overwrite an earlier SortType")
}

// Regression tests for https://github.com/Evertras/bubble-table/issues/232:
// SortTypeAuto used to pick numeric-vs-time-vs-string per pair of values
// rather than once for the whole column, which isn't transitive and made
// sort.Stable's result depend on row insertion order.

const mixedColumnKey = "v"

func sortColumnValues(t *testing.T, order []any) []any {
	t.Helper()

	rows := make([]Row, len(order))
	for i, v := range order {
		rows[i] = NewRow(RowData{mixedColumnKey: v})
	}

	model := New([]Column{
		NewColumn(mixedColumnKey, "V", 8),
	}).WithRows(rows).SortByAsc(mixedColumnKey)

	out := make([]any, len(order))

	for i, r := range model.GetVisibleRows() {
		out[i] = r.Data[mixedColumnKey]
	}

	return out
}

func TestSortMixedNumberStringColumnOrderIndependent(t *testing.T) {
	const valTenPlus = "10+"

	// A column mixing numbers and a string falls back to comparing every
	// pair as strings, so "10" < "10+" < "9" lexicographically ('1' < '9').
	expected := []any{10, valTenPlus, 9}

	assert.Equal(t, expected, sortColumnValues(t, []any{9, 10, valTenPlus}))
	assert.Equal(t, expected, sortColumnValues(t, []any{valTenPlus, 10, 9}))
	assert.Equal(t, expected, sortColumnValues(t, []any{10, 9, valTenPlus}))
	assert.Equal(t, expected, sortColumnValues(t, []any{9, valTenPlus, 10}))
}

func TestSortMixedTimeStringColumnOrderIndependent(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	// timeA is chronologically after timeB (15:00 UTC vs 20:00 UTC-1-day... below),
	// but the point is only that a mixed time/string column no longer cycles.
	timeA := time.Date(2021, 1, 1, 0, 0, 0, 0, jst)         // 2020-12-31 15:00 UTC
	timeB := time.Date(2020, 12, 31, 20, 0, 0, 0, time.UTC) // 2020-12-31 20:00 UTC
	strC := "mid"

	first := sortColumnValues(t, []any{timeA, timeB, strC})
	second := sortColumnValues(t, []any{strC, timeB, timeA})
	third := sortColumnValues(t, []any{timeB, strC, timeA})

	assert.Equal(t, first, second)
	assert.Equal(t, first, third)
}

func TestDetectAutoKind(t *testing.T) {
	const key = "v"

	makeRows := func(vals ...any) []Row {
		out := make([]Row, len(vals))
		for i, v := range vals {
			out[i] = NewRow(RowData{key: v})
		}

		return out
	}

	assert.Equal(t, autoKindNumeric, detectAutoKind(makeRows(1, 2.5, int64(3)), key))
	assert.Equal(t, autoKindTime, detectAutoKind(makeRows(time.Now(), time.Now()), key))
	assert.Equal(t, autoKindString, detectAutoKind(makeRows(1, "a"), key))
	assert.Equal(t, autoKindString, detectAutoKind(nil, key))
	assert.Equal(t, autoKindString, detectAutoKind([]Row{NewRow(RowData{})}, key))

	// A missing value in an otherwise-numeric column doesn't disqualify it.
	assert.Equal(t, autoKindNumeric, detectAutoKind([]Row{
		NewRow(RowData{key: 1}),
		NewRow(RowData{}),
	}, key))
}

func TestSortNumericColumnWithMissingValue(t *testing.T) {
	const idColKey = "id"

	rows := []Row{
		NewRow(RowData{idColKey: 5}),
		NewRow(RowData{}),
		NewRow(RowData{idColKey: 1}),
		NewRow(RowData{}),
	}

	model := New([]Column{NewColumn(idColKey, "ID", 3)}).WithRows(rows).SortByAsc(idColKey)

	_, existsAt0 := model.GetVisibleRows()[0].Data[idColKey]
	assert.False(t, existsAt0, "missing values should sort first ascending")
	_, existsAt1 := model.GetVisibleRows()[1].Data[idColKey]
	assert.False(t, existsAt1, "missing values should sort first ascending")
	assert.Equal(t, 1, model.GetVisibleRows()[2].Data[idColKey])
	assert.Equal(t, 5, model.GetVisibleRows()[3].Data[idColKey])

	model = model.SortByDesc(idColKey)

	assert.Equal(t, 5, model.GetVisibleRows()[0].Data[idColKey])
	assert.Equal(t, 1, model.GetVisibleRows()[1].Data[idColKey])
	_, existsAt2 := model.GetVisibleRows()[2].Data[idColKey]
	assert.False(t, existsAt2, "missing values should sort last descending")
	_, existsAt3 := model.GetVisibleRows()[3].Data[idColKey]
	assert.False(t, existsAt3, "missing values should sort last descending")
}

func TestSortTimeColumnWithMissingValue(t *testing.T) {
	const timeKey = "time"

	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	rows := []Row{
		NewRow(RowData{timeKey: base}),
		NewRow(RowData{}),
	}

	model := New([]Column{NewColumn(timeKey, "Time", 20)}).WithRows(rows).SortByAsc(timeKey)

	_, existsAfterAsc := model.GetVisibleRows()[0].Data[timeKey]
	assert.False(t, existsAfterAsc, "missing value should sort first ascending")

	model = model.SortByDesc(timeKey)

	_, existsAfterDesc := model.GetVisibleRows()[1].Data[timeKey]
	assert.False(t, existsAfterDesc, "missing value should sort last descending")
}

func TestSortTimeColumnTimezones(t *testing.T) {
	const timeKey = "time"

	// 8 AM EST is 1 PM UTC — later than 12 PM UTC
	utc := time.UTC
	est := time.FixedZone("EST", -5*60*60)

	noon := time.Date(2024, 1, 1, 12, 0, 0, 0, utc) // 12:00 UTC
	early := time.Date(2024, 1, 1, 8, 0, 0, 0, est) // 08:00 EST = 13:00 UTC (after noon)

	rows := []Row{
		NewRow(RowData{timeKey: early}), // chronologically later
		NewRow(RowData{timeKey: noon}),  // chronologically earlier
	}

	model := New([]Column{NewColumn(timeKey, "Time", 30)}).
		WithRows(rows).
		SortByAsc(timeKey)

	firstVal, ok := model.GetVisibleRows()[0].Data[timeKey].(time.Time)
	assert.True(t, ok)

	secondVal, ok := model.GetVisibleRows()[1].Data[timeKey].(time.Time)
	assert.True(t, ok)

	assert.True(t, firstVal.Before(secondVal) || firstVal.Equal(secondVal),
		"ascending: first row (%v) should be before second (%v)", firstVal, secondVal)
}
