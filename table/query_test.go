package table

import (
	"testing"

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
