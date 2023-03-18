package table

import (
	"fmt"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func TestIsRowMatchedForStyled(t *testing.T) {
	columns := []Column{
		NewColumn("title", "title", 10).WithFiltered(true),
		NewColumn("description", "description", 10),
	}

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}), "AA"))

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"title":       NewStyledCell("AAA", lipgloss.NewStyle()),
			"description": "",
		}), "AA"))
}

func TestIsRowMatchedForNonStringer(t *testing.T) {
	columns := []Column{
		NewColumn("val", "val", 10).WithFiltered(true),
	}

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"val": 12,
		}), "12"))

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"val": 12,
		}), "1"))

	assert.True(t, isRowMatched(columns,
		NewRow(RowData{
			"val": 12,
		}), "2"))

	assert.False(t, isRowMatched(columns,
		NewRow(RowData{
			"val": 12,
		}), "3"))
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

func TestGetFilteredRowsRefocusAfterFilter(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			"title":       "a",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "b",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "c",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "d1",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "d2",
			"description": "",
		}),
	}
	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(1)
	model = model.PageDown()
	assert.Len(t, model.GetVisibleRows(), 5)
	assert.Equal(t, 1, model.PageSize())
	assert.Equal(t, 2, model.CurrentPage())
	assert.Equal(t, 5, model.MaxPages())
	assert.Equal(t, 5, model.TotalRows())

	model.filterTextInput.SetValue("c")
	model, _ = model.updateFilterTextInput(tea.KeyMsg{})
	assert.Len(t, model.GetVisibleRows(), 1)
	assert.Equal(t, 1, model.PageSize())
	assert.Equal(t, 1, model.CurrentPage())
	assert.Equal(t, 1, model.MaxPages())
	assert.Equal(t, 1, model.TotalRows())

	model.filterTextInput.SetValue("not-exist")
	model, _ = model.updateFilterTextInput(tea.KeyMsg{})
	assert.Len(t, model.GetVisibleRows(), 0)
	assert.Equal(t, 1, model.PageSize())
	assert.Equal(t, 1, model.CurrentPage())
	assert.Equal(t, 1, model.MaxPages())
	assert.Equal(t, 0, model.TotalRows())
}

func TestFilterWithExternalTextInput(t *testing.T) {
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

	// Page size 1 to test scrolling back if input changes
	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(1)
	model.pageDown()
	assert.Equal(t, 2, model.CurrentPage(), "Should start on second page for test")
	input := textinput.New()
	input.SetValue("AaA")
	model = model.WithFilterInput(input)
	assert.Equal(t, 1, model.CurrentPage(), "Did not go back to first page")

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, 1)
}

func TestFilterWithSetValue(t *testing.T) {
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

	// Page size 1 to make sure we scroll back correctly
	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(1)
	model.pageDown()
	assert.Equal(t, 2, model.CurrentPage(), "Should start on second page for test")
	model = model.WithFilterInputValue("AaA")

	assert.Equal(t, 1, model.CurrentPage(), "Did not go back to first page")

	filteredRows := model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 1)

	// Make sure it holds true after an update
	model, _ = model.Update(tea.KeyRight)
	filteredRows = model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 1)

	// Remove filter
	model = model.WithFilterInputValue("")
	filteredRows = model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 3)
}

func BenchmarkFilteredScrolling(b *testing.B) {
	// Scrolling through a filtered table with many rows should be quick
	// https://github.com/Evertras/bubble-table/issues/135
	const rowCount = 40000
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := make([]Row, rowCount)

	for i := 0; i < rowCount; i++ {
		rows[i] = NewRow(RowData{
			"title": fmt.Sprintf("%d", i),
		})
	}

	model := New(columns).WithRows(rows).Filtered(true)
	model = model.WithFilterInputValue("1")

	hitKey := func(key rune) {
		model, _ = model.Update(
			tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{key},
			})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hitKey('j')
	}
}

func BenchmarkFilteredScrollingPaged(b *testing.B) {
	// Scrolling through a filtered table with many rows should be quick
	// https://github.com/Evertras/bubble-table/issues/135
	const rowCount = 40000
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := make([]Row, rowCount)

	for i := 0; i < rowCount; i++ {
		rows[i] = NewRow(RowData{
			"title": fmt.Sprintf("%d", i),
		})
	}

	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(50)
	model = model.WithFilterInputValue("1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model, _ = model.Update(
			tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'j'},
			})
	}
}

func BenchmarkFilteredRenders(b *testing.B) {
	// Rendering a filtered table should be fast
	// https://github.com/Evertras/bubble-table/issues/135
	const rowCount = 40000
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := make([]Row, rowCount)

	for i := 0; i < rowCount; i++ {
		rows[i] = NewRow(RowData{
			"title": fmt.Sprintf("%d", i),
		})
	}

	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(50)
	model = model.WithFilterInputValue("1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Don't care about result, just rendering
		_ = model.View()
	}
}
