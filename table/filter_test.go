package table

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestIsRowMatched(t *testing.T) {
	columns := []Column{
		NewColumn("title", "title", 10).WithFiltered(true),
		NewColumn("description", "description", 10)}

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), ""))

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "AA"))

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "A"))

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "a"))

	assert.False(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "B"))

	assert.False(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "BBB",
		}), "BBB"))
}

func TestGetFilteredRows(t *testing.T) {
	input := textinput.Model{}
	input.SetValue("AAA")
	columns := []Column{NewColumn("title", "title", 10)}
	model := Model{filtered: true, filterTextInput: input, columns: columns}
	rows := []Row{
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(map[string]interface{}{
			"title":       "BBB",
			"description": "",
		}),
		NewRow(map[string]interface{}{
			"title":       "CCC",
			"description": "",
		}),
	}
	filteredRows := model.getFilteredRows(rows)
	assert.Equal(t, 0, len(filteredRows))
}

func TestGetFilteredRowsFiltered(t *testing.T) {
	input := textinput.Model{}
	input.SetValue("AAA")
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	model := Model{filtered: true, filterTextInput: input, columns: columns}
	rows := []Row{
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(map[string]interface{}{
			"title":       "BBB",
			"description": "",
		}),
		NewRow(map[string]interface{}{
			"title":       "CCC",
			"description": "",
		}),
	}
	filteredRows := model.getFilteredRows(rows)
	assert.Equal(t, 1, len(filteredRows))
}
