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
	})
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
