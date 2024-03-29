package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithHighlightedRowSet(t *testing.T) {
	highlightedIndex := 1

	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
	}).WithHighlightedRow(highlightedIndex)

	assert.Equal(t, model.rows[highlightedIndex], model.HighlightedRow())
}

func TestWithHighlightedRowSetNegative(t *testing.T) {
	highlightedIndex := -1

	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
	}).WithHighlightedRow(highlightedIndex)

	assert.Equal(t, model.rows[0], model.HighlightedRow())
}

func TestWithHighlightedRowSetTooHigh(t *testing.T) {
	highlightedIndex := 2

	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
	}).WithHighlightedRow(highlightedIndex)

	assert.Equal(t, model.rows[1], model.HighlightedRow())
}

// This is long only because it's a lot of repetitive test cases
//
//nolint:funlen
func TestPageOptions(t *testing.T) {
	const (
		pageSize = 5
		rowCount = 30
	)

	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	rows := make([]Row, rowCount)

	model := New(cols).WithRows(rows).WithPageSize(pageSize)
	assert.Equal(t, 1, model.CurrentPage())

	model = model.PageDown()
	assert.Equal(t, 2, model.CurrentPage())

	model = model.PageDown()
	model = model.PageUp()
	assert.Equal(t, 2, model.CurrentPage())

	model = model.PageLast()
	assert.Equal(t, 6, model.CurrentPage())

	model = model.PageLast()
	model = model.PageLast()
	assert.Equal(t, 6, model.CurrentPage())

	model = model.PageFirst()
	assert.Equal(t, 1, model.CurrentPage())

	model = model.PageFirst()
	model = model.PageFirst()
	assert.Equal(t, 1, model.CurrentPage())

	model = model.PageUp()
	assert.Equal(t, 6, model.CurrentPage())

	model = model.PageDown()
	assert.Equal(t, 1, model.CurrentPage())

	model = model.WithCurrentPage(3)
	model = model.WithCurrentPage(3)
	model = model.WithCurrentPage(3)
	assert.Equal(t, 3, model.CurrentPage())
	assert.Equal(t, 10, model.rowCursorIndex)

	model = model.WithCurrentPage(-1)
	assert.Equal(t, 1, model.CurrentPage())
	assert.Equal(t, 0, model.rowCursorIndex)

	model = model.WithCurrentPage(0)
	assert.Equal(t, 1, model.CurrentPage())
	assert.Equal(t, 0, model.rowCursorIndex)

	model = model.WithCurrentPage(7)
	assert.Equal(t, 6, model.CurrentPage())
	assert.Equal(t, 25, model.rowCursorIndex)

	model.rowCursorIndex = 26
	model = model.WithCurrentPage(6)
	assert.Equal(t, 6, model.CurrentPage())
	assert.Equal(t, 26, model.rowCursorIndex)

	model = model.WithFooterVisibility(false)
	assert.Equal(t, "", model.renderFooter(10, false))

	model = model.WithFooterVisibility(true)
	assert.Greater(t, len(model.renderFooter(10, false)), 10)
	assert.Contains(t, model.renderFooter(10, false), "6/6")
}

func TestMinimumHeightOptions(t *testing.T) {
	columns := []Column{
		NewColumn("ka", "a", 3),
		NewColumn("kb", "b", 4),
		NewColumn("kc", "c", 5),
	}

	model := New(columns).WithMinimumHeight(10)
	assert.Equal(t, 10, model.minimumHeight)
	assert.Equal(t, 3, model.metaHeight)

	model = model.WithPageSize(2)
	assert.Equal(t, 5, model.metaHeight)

	model = model.WithNoPagination()
	assert.Equal(t, 3, model.metaHeight)

	model = model.WithStaticFooter("footer with\nmultiple lines")
	assert.Equal(t, 6, model.metaHeight)

	model = model.WithStaticFooter("").Filtered(true)
	assert.Equal(t, 5, model.metaHeight)

	model = model.WithFooterVisibility(false)
	assert.Equal(t, 3, model.metaHeight)

	model = model.WithHeaderVisibility(false)
	assert.Equal(t, 1, model.metaHeight)
}

// This is long only because the test cases are larger
//
//nolint:funlen
func TestSelectRowsProgramatically(t *testing.T) {
	const col = "id"

	tests := map[string]struct {
		rows        []Row
		selectedIds []int
	}{
		"no rows selected": {
			[]Row{
				NewRow(RowData{col: 1}),
				NewRow(RowData{col: 2}),
				NewRow(RowData{col: 3}),
			},
			[]int{},
		},

		"all rows selected": {
			[]Row{
				NewRow(RowData{col: 1}).Selected(true),
				NewRow(RowData{col: 2}).Selected(true),
				NewRow(RowData{col: 3}).Selected(true),
			},
			[]int{1, 2, 3},
		},

		"first row selected": {
			[]Row{
				NewRow(RowData{col: 1}).Selected(true),
				NewRow(RowData{col: 2}),
				NewRow(RowData{col: 3}),
			},
			[]int{1},
		},

		"last row selected": {
			[]Row{
				NewRow(RowData{col: 1}),
				NewRow(RowData{col: 2}),
				NewRow(RowData{col: 3}).Selected(true),
			},
			[]int{3},
		},
	}

	baseModel := New([]Column{
		NewColumn(col, col, 1),
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			model := baseModel.WithRows(test.rows)
			sel := model.SelectedRows()

			assert.Equal(t, len(test.selectedIds), len(sel))
			for i, id := range test.selectedIds {
				assert.Equal(t, id, sel[i].Data[col], "expecting row %d to have same %s column value", i)
			}

			model = model.WithAllRowsDeselected()
			assert.Len(t, model.SelectedRows(), 0, "Did not deselect all rows")
		})
	}
}

func BenchmarkSelectedRows(b *testing.B) {
	const N = 1000

	b.ReportAllocs()

	rows := make([]Row, 0, N)
	for i := 0; i < N; i++ {
		rows = append(rows, NewRow(RowData{"row": i}).Selected(i%2 == 0))
	}

	model := New([]Column{
		NewColumn("row", "Row", 4),
	}).WithRows(rows)

	var sel []Row

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sel = model.SelectedRows()
	}

	Rows = sel
}

var Rows []Row
