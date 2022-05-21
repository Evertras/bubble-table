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
// nolint: funlen
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
	assert.Equal(t, "", model.renderFooter())

	model = model.WithFooterVisibility(true)
	assert.Greater(t, len(model.renderFooter()), 10)
	assert.Contains(t, model.renderFooter(), "6/6")
}
