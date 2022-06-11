package table

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestGetColumnSorting(t *testing.T) {
	cols := []Column{
		NewColumn("a", "a", 3),
		NewColumn("b", "b", 3),
		NewColumn("c", "c", 3),
	}

	model := New(cols).SortByAsc("b")

	sorted := model.GetColumnSorting()

	assert.Len(t, sorted, 1, "Should only have one column")
	assert.Equal(t, sorted[0].ColumnKey, "b", "Should sort column b")
	assert.Equal(t, sorted[0].Direction, SortDirectionAsc, "Should be ascending")

	sorted[0].Direction = SortDirectionDesc

	assert.NotEqual(
		t,
		model.sortOrder[0].Direction,
		sorted[0].Direction,
		"Should not have been able to modify actual values",
	)
}

func TestGetFilterData(t *testing.T) {
	model := New([]Column{})

	assert.False(t, model.GetIsFilterActive(), "Should not start with filter active")
	assert.False(t, model.GetCanFilter(), "Should not start with filter ability")
	assert.Equal(t, model.GetCurrentFilter(), "", "Filter string should be empty")

	model = model.Filtered(true)

	assert.False(t, model.GetIsFilterActive(), "Should not be filtered just because the ability was activated")
	assert.True(t, model.GetCanFilter(), "Filter feature should be enabled")
	assert.Equal(t, model.GetCurrentFilter(), "", "Filter string should be empty")

	model.filterTextInput.SetValue("a")

	assert.True(t, model.GetIsFilterActive(), "Typing anything into box should mark as filtered")
	assert.True(t, model.GetCanFilter(), "Filter feature should be enabled")
	assert.Equal(t, model.GetCurrentFilter(), "a", "Filter string should be what was typed")
}

func TestGetVisibleRows(t *testing.T) {
	input := textinput.Model{}
	input.SetValue("AAA")
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
		NewRow(RowData{
			"title":       "CCC",
			"description": "",
		}),
	}
	m := Model{filtered: true, filterTextInput: input, columns: columns, rows: rows}
	visibleRows := m.GetVisibleRows()
	assert.Len(t, visibleRows, 1)
}

func TestGetHighlightedRowIndex(t *testing.T) {
	model := New([]Column{})

	assert.Equal(t, 0, model.GetHighlightedRowIndex(), "Empty table should still safely have 0 index highlighted")

	// We don't actually need data to test this
	empty := RowData{}
	model = model.WithRows([]Row{NewRow(empty), NewRow(empty)})

	assert.Equal(t, 0, model.GetHighlightedRowIndex(), "Unfocused table should start with 0 index")

	model = model.WithHighlightedRow(1)

	assert.Equal(t, 1, model.GetHighlightedRowIndex(), "Table with set highlighted row should return same highlighted row")
}

func TestGetFocused(t *testing.T) {
	model := New([]Column{})

	assert.Equal(t, false, model.GetFocused(), "Table should not be focused by default")

	model = model.Focused(true)

	assert.Equal(t, true, model.GetFocused(), "Table should be focused after being set")
}

func TestGetHorizontalScrollColumnOffset(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
		NewColumn("4", "4", 4),
	}).
		WithRows([]Row{
			NewRow(RowData{
				"1": "x1",
				"2": "x2",
				"3": "x3",
				"4": "x4",
			}),
		}).
		WithMaxTotalWidth(18).
		Focused(true)

	hitScrollRight := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
	}

	hitScrollLeft := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftLeft})
	}

	assert.Equal(
		t,
		0,
		model.GetHorizontalScrollColumnOffset(),
		"Should start to left",
	)

	hitScrollRight()

	assert.Equal(
		t,
		1,
		model.GetHorizontalScrollColumnOffset(),
		"Should be 1 after scrolling to the right once",
	)

	hitScrollLeft()
	assert.Equal(
		t,
		0,
		model.GetHorizontalScrollColumnOffset(),
		"Should be back to 0 after moving to the left",
	)

	hitScrollLeft()
	assert.Equal(
		t,
		0,
		model.GetHorizontalScrollColumnOffset(),
		"Should still be 0 after trying to go left again",
	)
}

func TestGetHeaderVisibility(t *testing.T) {
	model := New([]Column{})

	assert.True(t, model.GetHeaderVisibility(), "Header should be visible by default")

	model = model.WithHeaderVisibility(false)

	assert.False(t, model.GetHeaderVisibility(), "Header was not set to hidden")
}
