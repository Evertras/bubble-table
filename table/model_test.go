package table

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestModelInitReturnsNil(t *testing.T) {
	model := New(nil)

	cmd := model.Init()

	assert.Nil(t, cmd)
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
