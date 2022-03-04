package table

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsRowMatched(t *testing.T) {
	columns := []Column{
		NewColumn("title", "title", 10).WithFiltered(true),
		NewColumn("description", "description", 10)}

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}), ""))

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}), "AA"))

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}), "A"))

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}), "a"))

	assert.False(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}), "B"))

	assert.False(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "BBB",
		}), "BBB"))

	timeFrom2020 := time.Date(2020, time.July, 1, 1, 1, 1, 1, time.UTC)

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title": timeFrom2020,
		}),
		"2020",
	))

	assert.False(t, isRowMatched(columns,
		NewRow(RowData{
			"title": timeFrom2020,
		}),
		"2021",
	))
}

func TestGetFilteredRowsNoColumnFiltered(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10)}
	rows := []Row{
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "BBB",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "CCC",
			"description": "",
		}),
	}

	model := New(columns).WithRows(rows).Filtered(true)
	model.filterTextInput.SetValue("AAA")

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, len(rows))
}

func TestGetFilteredRowsUnfiltered(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10)}
	rows := []Row{
		NewRow(RowData{
			"title": "AAA",
		}),
		NewRow(RowData{
			"title": "BBB",
		}),
	}

	model := New(columns).WithRows(rows)

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, len(rows))
}

func TestGetFilteredRowsFiltered(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "BBB",
			"description": "",
		}),
		// Empty
		NewRow(RowData{}),
	}
	model := New(columns).WithRows(rows).Filtered(true)
	model.filterTextInput.SetValue("AaA")

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, 1)
}
