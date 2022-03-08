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
