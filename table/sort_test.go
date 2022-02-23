package table

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestSortSingleColumnAsc(t *testing.T) {
	const idColKey = "id"
	model := New([]Column{
		NewColumn(idColKey, "ID", 3),
	}).WithRows([]Row{
		NewRow(RowData{idColKey: "b"}),
		NewRow(RowData{idColKey: NewStyledCell("c", lipgloss.NewStyle().Bold(true))}),
		NewRow(RowData{idColKey: "a"}),
	}).SortByAsc(idColKey)

	colKeyForRow := func(i int) string {
		id, ok := model.sortedRows[i].Data[idColKey]

		assert.True(t, ok)

		switch id := id.(type) {
		case string:
			return id

		case StyledCell:
			return id.Data.(string)

		default:
			assert.FailNowf(t, "Unknown type", "Unknown type in %d", i)
			return ""
		}
	}

	assert.Len(t, model.sortedRows, 3)
	assert.Equal(t, "a", colKeyForRow(0))
	assert.Equal(t, "b", colKeyForRow(1))
	assert.Equal(t, "c", colKeyForRow(2))
}

func TestSortSingleColumnDesc(t *testing.T) {
	const idColKey = "id"
	model := New([]Column{
		NewColumn(idColKey, "ID", 3),
	}).WithRows([]Row{
		NewRow(RowData{idColKey: "b"}),
		NewRow(RowData{idColKey: NewStyledCell("c", lipgloss.NewStyle().Bold(true))}),
		NewRow(RowData{idColKey: "a"}),
	}).SortByDesc(idColKey)

	colKeyForRow := func(i int) string {
		id, ok := model.sortedRows[i].Data[idColKey]

		assert.True(t, ok)

		switch id := id.(type) {
		case string:
			return id

		case StyledCell:
			return id.Data.(string)

		default:
			assert.FailNowf(t, "Unknown type", "Unknown type in %d", i)
			return ""
		}
	}

	assert.Len(t, model.sortedRows, 3)
	assert.Equal(t, "c", colKeyForRow(0))
	assert.Equal(t, "b", colKeyForRow(1))
	assert.Equal(t, "a", colKeyForRow(2))
}

func TestSortSingleColumnMissingValues(t *testing.T) {
	const idColKey = "id"
	model := New([]Column{
		NewColumn(idColKey, "ID", 3),
	}).WithRows([]Row{
		NewRow(RowData{idColKey: "b"}),
		NewRow(RowData{idColKey: NewStyledCell("c", lipgloss.NewStyle().Bold(true))}),
		NewRow(RowData{}),
	}).SortByAsc(idColKey)

	colKeyForRow := func(i int) string {
		id, ok := model.sortedRows[i].Data[idColKey]

		if !ok {
			id = ""
		}

		switch id := id.(type) {
		case string:
			return id

		case StyledCell:
			return id.Data.(string)

		default:
			return ""
		}
	}

	assert.Len(t, model.sortedRows, 3)
	assert.Equal(t, "", colKeyForRow(0))
	assert.Equal(t, "b", colKeyForRow(1))
	assert.Equal(t, "c", colKeyForRow(2))
}

func TestSortSingleColumnIntsAsc(t *testing.T) {
	const idColKey = "id"
	model := New([]Column{
		NewColumn(idColKey, "ID", 3),
	}).WithRows([]Row{
		NewRow(RowData{idColKey: 13}),
		NewRow(RowData{idColKey: NewStyledCell(1, lipgloss.NewStyle().Bold(true))}),
		NewRow(RowData{idColKey: 2}),
	}).SortByAsc(idColKey)

	colKeyForRow := func(i int) int {
		id, ok := model.sortedRows[i].Data[idColKey]

		assert.True(t, ok)

		switch id := id.(type) {
		case int:
			return id

		case StyledCell:
			return id.Data.(int)

		default:
			assert.FailNowf(t, "Unknown type", "Unknown type in %d", i)
			return 0
		}
	}

	assert.Len(t, model.sortedRows, 3)
	assert.Equal(t, 1, colKeyForRow(0))
	assert.Equal(t, 2, colKeyForRow(1))
	assert.Equal(t, 13, colKeyForRow(2))
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

	assertVals := func(i int, name string, score int) {
		actualName, ok := model.sortedRows[i].Data[nameKey].(string)
		assert.True(t, ok)

		actualScore, ok := model.sortedRows[i].Data[scoreKey].(int)
		assert.True(t, ok)

		assert.Equal(t, name, actualName)
		assert.Equal(t, score, actualScore)
	}

	assert.Len(t, model.sortedRows, 4)

	assertVals(0, "a", 100)
	assertVals(1, "a", 75)
	assertVals(2, "b", 101)
	assertVals(3, "c", 50)
}

func TestSortTwoColumnsDescAscMix(t *testing.T) {
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
	}).SortByDesc(nameKey).ThenSortByAsc(scoreKey)

	assertVals := func(i int, name string, score int) {
		actualName, ok := model.sortedRows[i].Data[nameKey].(string)
		assert.True(t, ok)

		actualScore, ok := model.sortedRows[i].Data[scoreKey].(int)
		assert.True(t, ok)

		assert.Equal(t, name, actualName)
		assert.Equal(t, score, actualScore)
	}

	assert.Len(t, model.sortedRows, 4)

	assertVals(0, "c", 50)
	assertVals(1, "b", 101)
	assertVals(2, "a", 75)
	assertVals(3, "a", 100)
}
