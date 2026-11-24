package table

import (
	"fmt"
	"regexp"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

const kindKey = "kind"

var ctrlF = tea.KeyPressMsg{Code: 'f', Mod: tea.ModCtrl}

func genSearchTable() Model {
	return New([]Column{
		NewColumn("name", "Name", 10),
		NewColumn(kindKey, "Kind", 10),
	}).WithRows([]Row{
		NewRow(RowData{"name": "apple", kindKey: "fruit"}),
		NewRow(RowData{"name": "carrot", kindKey: "veg"}),
		NewRow(RowData{"name": "banana", kindKey: "fruit"}),
		NewRow(RowData{"name": "leek", kindKey: "veg"}),
	}).Focused(true)
}

func typeKeys(model Model, text string) Model {
	for _, r := range text {
		model, _ = model.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}

	return model
}

func TestSearchNextPrevWrapAround(t *testing.T) {
	// Default search covers all columns, even ones not marked filterable
	model := genSearchTable().WithSearchInputValue("FRUIT")

	model = model.SearchNext()
	assert.Equal(t, 2, model.rowCursorIndex)

	model = model.SearchNext()
	assert.Equal(t, 0, model.rowCursorIndex)

	model = model.SearchPrev()
	assert.Equal(t, 2, model.rowCursorIndex)
}

func TestSearchNoMatchOrEmptyQueryKeepsCursor(t *testing.T) {
	model := genSearchTable().WithHighlightedRow(1)

	assert.Equal(t, 1, model.SearchNext().rowCursorIndex)
	assert.Equal(t, 1, model.WithSearchInputValue("zzz").SearchNext().rowCursorIndex)
	assert.Equal(t, 1, model.WithSearchInputValue("zzz").SearchPrev().rowCursorIndex)
	assert.Equal(t, 0, New(nil).WithSearchInputValue("x").SearchNext().rowCursorIndex)
}

func TestSearchRespectsFilter(t *testing.T) {
	model := New([]Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
		NewColumn(kindKey, "Kind", 10),
	}).WithRows(genSearchTable().rows).Filtered(true).
		WithFilterInputValue("a"). // apple, carrot, banana
		WithSearchInputValue("veg")

	model = model.SearchNext()
	assert.Equal(t, 1, model.rowCursorIndex)
	assert.Equal(t, "carrot", model.HighlightedRow().Data["name"])
}

func TestSearchCustomFunc(t *testing.T) {
	model := genSearchTable().
		WithSearchFunc(func(input FilterFuncInput) bool {
			ok, _ := regexp.MatchString(input.Filter, fmt.Sprint(input.Row.Data["name"]))

			return ok
		}).
		WithSearchInputValue("^l")

	assert.Equal(t, 3, model.SearchNext().rowCursorIndex)
	assert.Equal(t, 3, model.WithSearchFunc(nil).WithSearchInputValue("leek").SearchNext().rowCursorIndex)
}

func TestSearchFuzzy(t *testing.T) {
	model := genSearchTable().
		WithFuzzySearch().
		WithSearchInputValue("crrt") // subsequence of "carrot"

	assert.Equal(t, 1, model.SearchNext().rowCursorIndex)
}

func TestSearchChangesPage(t *testing.T) {
	model := genSearchTable().WithPageSize(2).WithSearchInputValue("leek").SearchNext()

	assert.Equal(t, 3, model.rowCursorIndex)
	assert.Equal(t, 2, model.CurrentPage())
}

func TestSearchKeysIncremental(t *testing.T) {
	model := genSearchTable().Searchable(true)

	model, _ = model.Update(ctrlF)
	assert.True(t, model.GetIsSearchInputFocused())

	// Incremental: jumps as the query is typed, keys go to the input not the table
	model = typeKeys(model, "b")
	assert.Equal(t, 2, model.rowCursorIndex)
	assert.Contains(t, model.GetLastUpdateUserEvents(), UserEventHighlightedIndexChanged{
		PreviousRowIndex: 0,
		SelectedRowIndex: 2,
	})

	// Deleting the query goes back to where the search started
	model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	assert.Equal(t, 0, model.rowCursorIndex)

	model = typeKeys(model, "ba")
	assert.Equal(t, 2, model.rowCursorIndex)
	assert.Equal(t, "ba", model.GetCurrentSearch())
	assert.Contains(t, model.View(), "?ba")

	// Keys that don't change the query (moving the text cursor) don't search again
	model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	assert.Equal(t, 2, model.rowCursorIndex)
	assert.Empty(t, model.GetLastUpdateUserEvents())

	model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.False(t, model.GetIsSearchInputFocused())
	assert.Contains(t, model.View(), "?ba")

	model = model.WithSearchInputValue("veg")
	model = typeKeys(model, "n")
	assert.Equal(t, 3, model.rowCursorIndex)

	model = typeKeys(model, "N")
	assert.Equal(t, 1, model.rowCursorIndex)

	// Esc while blurred clears the search along with the filter
	model, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	assert.Empty(t, model.GetCurrentSearch())
}

func TestSearchKeysIgnoredWhenNotSearchable(t *testing.T) {
	model, _ := genSearchTable().WithSearchInputValue("leek").Update(ctrlF)
	model = typeKeys(model, "n")

	assert.False(t, model.GetIsSearchInputFocused())
	assert.Equal(t, 0, model.rowCursorIndex)
}

func TestRowStyleFuncIsSearchMatch(t *testing.T) {
	var matched []int

	model := genSearchTable().
		WithSearchInputValue("veg").
		WithRowStyleFunc(func(input RowStyleFuncInput) lipgloss.Style {
			if input.IsSearchMatch {
				matched = append(matched, input.Index)
			}

			return lipgloss.NewStyle()
		})

	model.View()

	assert.Equal(t, []int{1, 3}, matched)
}

func TestSearchableTargetHeightAndHelp(t *testing.T) {
	model := genSearchTable().WithTargetHeight(6).WithSearchInputValue("leek")

	assert.False(t, model.GetCanSearch())
	assert.Len(t, model.FullHelp(), 3)

	model = model.Searchable(true)

	assert.True(t, model.GetCanSearch())
	assert.Len(t, model.FullHelp(), 4)
	assert.Equal(t, 3, model.SearchPrev().rowCursorIndex)
	assert.Contains(t, model.View(), "?leek")
}
