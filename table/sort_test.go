package table

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestSortSingleColumnAscAndDesc(t *testing.T) {
	const idColKey = "id"

	rows := []Row{
		NewRow(RowData{idColKey: "b"}),
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
	sortColumns := []sortColumn{
		{
			columnKey: "cb",
			direction: 1,
		},
		{
			columnKey: "ca",
			direction: 0,
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
